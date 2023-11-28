package api

import "context"

func (r *galleryFileResolver) Fingerprint(ctx context.Context, obj *GalleryFile, type_ string) (*string, error) {
	fp := obj.BaseFile.Fingerprints.For(type_)
	if fp != nil {
		v := fp.Value()
		return &v, nil
	}
	return nil, nil
}

func (r *imageFileResolver) Fingerprint(ctx context.Context, obj *ImageFile, type_ string) (*string, error) {
	fp := obj.ImageFile.Fingerprints.For(type_)
	if fp != nil {
		v := fp.Value()
		return &v, nil
	}
	return nil, nil
}

func (r *videoFileResolver) Fingerprint(ctx context.Context, obj *VideoFile, type_ string) (*string, error) {
	fp := obj.VideoFile.Fingerprints.For(type_)
	if fp != nil {
		v := fp.Value()
		return &v, nil
	}
	return nil, nil
}
