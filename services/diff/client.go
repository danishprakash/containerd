package diff

import (
	"context"

	"github.com/containerd/containerd"
	diffapi "github.com/containerd/containerd/api/services/diff"
	"github.com/containerd/containerd/api/types/descriptor"
	mounttypes "github.com/containerd/containerd/api/types/mount"
	"github.com/containerd/containerd/rootfs"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// NewApplierFromClient returns a new Applier which communicates
// over a GRPC connection.
func NewApplierFromClient(client diffapi.DiffClient) rootfs.Applier {
	return &remoteApplier{
		client: client,
	}
}

type remoteApplier struct {
	client diffapi.DiffClient
}

func (r *remoteApplier) Apply(ctx context.Context, diff ocispec.Descriptor, mounts []containerd.Mount) (ocispec.Descriptor, error) {
	req := &diffapi.ApplyRequest{
		Diff:   fromDescriptor(diff),
		Mounts: fromMounts(mounts),
	}
	resp, err := r.client.Apply(ctx, req)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	return toDescriptor(resp.Applied), nil
}

func fromDescriptor(d ocispec.Descriptor) *descriptor.Descriptor {
	return &descriptor.Descriptor{
		MediaType: d.MediaType,
		Digest:    d.Digest,
		Size_:     d.Size,
	}
}

func fromMounts(mounts []containerd.Mount) []*mounttypes.Mount {
	apiMounts := make([]*mounttypes.Mount, len(mounts))
	for i, m := range mounts {
		apiMounts[i] = &mounttypes.Mount{
			Type:    m.Type,
			Source:  m.Source,
			Options: m.Options,
		}
	}
	return apiMounts
}
