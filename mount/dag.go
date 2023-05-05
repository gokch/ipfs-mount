package mount

import (
	"context"
	"io"

	"github.com/gokch/kioskgo/file"
	"github.com/ipfs/boxo/blockservice"
	"github.com/ipfs/boxo/exchange"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/ipld/merkledag"
	unixfile "github.com/ipfs/boxo/ipld/unixfs/file"
	"github.com/ipfs/boxo/ipld/unixfs/importer/balanced"
	uih "github.com/ipfs/boxo/ipld/unixfs/importer/helpers"
	"github.com/ipfs/go-cid"
	chunk "github.com/ipfs/go-ipfs-chunker"
	"github.com/multiformats/go-multicodec"
)

// Dag dag to fileStore
// TODO : fs 의 cid 와 dag 의 cid 가 다를 경우 동기화 처리 필요
// dag 의 block 은 어느 기준으로 Garbage collect?? 블록 전체를 캐싱하고 있으면 안되는데...
type Dag struct {
	blockSize int

	Dag   *uih.DagBuilderParams // use MapDataStore
	mount *Mount                // FileStore
}

func NewDag(ctx context.Context, blockSize int, mount *Mount, rem exchange.Interface) (*Dag, error) {
	// set default block size
	if blockSize <= 0 {
		blockSize = int(chunk.DefaultBlockSize)
	}
	// make dag service, save dht blocks
	// Create a UnixFS graph from our file, parameters described here but can be visualized at https://dag.ipfs.tech/
	builder := &uih.DagBuilderParams{
		Maxlinks:  uih.DefaultLinksPerBlock, // Default max of 174 links per block
		RawLeaves: true,                     // Leave the actual file bytes untouched instead of wrapping them in a dag-pb protobuf wrapper
		CidBuilder: cid.V1Builder{ // Use CIDv1 for all links
			Codec:    uint64(multicodec.DagPb),
			MhType:   uint64(multicodec.Sha3_256), // Use SHA3-256 as the hash function
			MhLength: -1,                          // Use the default hash length for the given hash function (in this case 256 bits)
		},
		Dagserv: merkledag.NewDAGService(blockservice.New(mount, rem)),
		NoCopy:  false,
	}

	Dag := &Dag{
		blockSize: blockSize,
		Dag:       builder,
		mount:     mount,
	}

	// import blocks in Merkle-DAG from fileStore
	err := Dag.mount.fs.Iterate("", func(path string, reader *file.Reader) error {
		_, err := Dag.Upload(ctx, path, nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return Dag, nil
}

func (d *Dag) Download(ctx context.Context, ci cid.Cid, path string) error {
	nd, err := d.Dag.Dagserv.Get(ctx, ci)
	if err != nil {
		return err
	}

	// put cids in fm
	d.mount.fm.PutNode(nd, path, d.blockSize)

	unixFSNode, err := unixfile.NewUnixfsFile(ctx, d.Dag.Dagserv, nd)
	if err != nil {
		return err
	}
	// put data in fs
	err = d.mount.fs.Put(ctx, path, file.NewWriter(unixFSNode))
	if err != nil {
		return err
	}
	unixFSNode.Close()

	return nil
}

func (d *Dag) Upload(ctx context.Context, path string, reader io.Reader) (cid.Cid, error) {
	var err error

	// if reader != nil, put reader in fileStore
	// if already exist in filepath, return err
	if reader != nil {
		err = d.mount.fs.Put(ctx, path, file.NewWriter(files.NewReaderFile(reader)))
		if err != nil {
			return cid.Cid{}, err
		}
	}

	// get reader from filestore
	reader, err = d.mount.fs.Get(ctx, path)
	if err != nil {
		return cid.Cid{}, err
	}

	// put reader in dag
	ufsBuilder, err := d.Dag.New(chunk.NewSizeSplitter(reader, 1000))
	if err != nil {
		return cid.Cid{}, err
	}
	nd, err := balanced.Layout(ufsBuilder) // Arrange the graph with a balanced layout
	if err != nil {
		return cid.Cid{}, err
	}

	// put cids in fm
	d.mount.fm.PutNode(nd, path, d.blockSize)

	return nd.Cid(), nil
}
