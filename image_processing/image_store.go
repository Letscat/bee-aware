package image_processing

type CameraStore struct {
	frames map[string][][]uint8
}

func NewCameraStore() *CameraStore {
	return &CameraStore{frames: make(map[string][][]uint8)}
}

func (cs *CameraStore) StoreFrame(cameraID string, pixels [][]uint8) {
	cs.frames[cameraID] = pixels
}

func (cs *CameraStore) GetLastFrame(cameraID string) ([][]uint8, bool) {
	pixels, ok := cs.frames[cameraID]
	return pixels, ok
}
