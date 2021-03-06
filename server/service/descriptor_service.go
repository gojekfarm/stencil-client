package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"path"
	"time"

	"github.com/odpf/stencil/server/models"
	"github.com/odpf/stencil/server/proto"
	"github.com/odpf/stencil/server/store"
)

//DescriptorService Interacts with backend store
type DescriptorService struct {
	Store *store.Store
}

//ListNames returns list of directories
func (d *DescriptorService) ListNames(prefixes ...string) ([]string, error) {
	prefix := path.Join(prefixes...)
	return d.Store.ListDir(prefix + "/")
}

//ListVersions returns list of versions for specified prefixes
func (d *DescriptorService) ListVersions(prefixes ...string) ([]string, error) {
	prefix := path.Join(prefixes...)
	return d.Store.ListFiles(prefix + "/")
}

//Upload uploads the file
func (d *DescriptorService) Upload(ctx context.Context, payload *models.DescriptorPayload) error {
	namespace, name, version := payload.Namespace, payload.Name, payload.Version
	filename := path.Join(namespace, name, version)
	exists, _ := d.Store.Exists(ctx, filename)
	if exists {
		return models.ErrConflict
	}
	data, err := readDataFromMultiPartFile(payload.File)
	if err != nil {
		return models.WrapAPIError(models.ErrUploadInvalidFile, err)
	}
	err = d.isBackwardCompatible(ctx, payload, data)
	if err != nil {
		return err
	}
	if payload.DryRun {
		return nil
	}
	err = d.Store.PutData(ctx, filename, data)
	if err != nil {
		return err
	}
	if payload.Latest {
		return d.StoreMetadata(ctx, &models.MetadataPayload{Version: version, Name: name, Namespace: namespace})
	}
	return nil
}

//Download downloads the file
func (d *DescriptorService) Download(ctx context.Context, payload *models.FileDownload) (*models.FileData, error) {
	filename := path.Join(payload.Namespace, payload.Name, payload.Version)
	data, err := d.Store.Get(ctx, filename)
	if err != nil {
		return nil, err
	}
	return &models.FileData{
		ContentLength: data.Size(),
		Reader:        data,
	}, nil
}

//StoreMetadata stores latest version number
func (d *DescriptorService) StoreMetadata(ctx context.Context, payload *models.MetadataPayload) error {
	prefix := path.Join(payload.Namespace, payload.Name)
	metafile := path.Join(prefix, "meta.json")
	filename := path.Join(prefix, payload.Version)
	fileExists, err := d.Store.Exists(ctx, filename)
	if !fileExists {
		return models.WrapAPIError(models.ErrNotFound, err)
	}
	updated := time.Now().UTC().Format(time.RFC3339)
	fileData := &models.MetadataFile{
		Version: payload.Version,
		Updated: updated,
	}
	data, err := json.Marshal(fileData)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(data)
	err = d.Store.Put(ctx, metafile, reader)
	if err != nil {
		return err
	}
	return d.Store.Copy(ctx, filename, path.Join(prefix, "latest"))
}

//GetMetadata gets latest version number
func (d *DescriptorService) GetMetadata(ctx context.Context, payload *models.GetMetadata) (*models.MetadataFile, error) {
	filename := path.Join(payload.Namespace, payload.Name, "meta.json")
	data, err := d.Store.Get(ctx, filename)
	if err != nil {
		return nil, err
	}
	defer data.Close()
	file := &models.MetadataFile{}
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (d *DescriptorService) isBackwardCompatible(ctx context.Context, payload *models.DescriptorPayload, data []byte) error {
	filename := path.Join(payload.Namespace, payload.Name, "latest")
	reader, err := d.Store.Get(ctx, filename)
	if err != nil {
		if isNotFoundErr(err) {
			return nil
		}
		return err
	}
	prevData, err := readDataFromReader(reader)
	if err != nil {
		return err
	}
	err = proto.Compare(data, prevData, payload.SkipRules)
	if err != nil {
		return models.NewAPIError(400, err.Error(), err)
	}
	return err
}
