// Service site keeps track of which sites to monitor.
package motion_detection

import (
	"context"

	"encore.app/image_processing"
	"encore.dev/beta/errs"
)

// Site describes a monitored site.
type Image struct {
	// ID is a unique ID for the site.
	FileData string `json:"fileData"`
}
type Motion struct {
	Movement bool `json:"movement"`
}

// AddParams are the parameters for adding a site to be monitored.
type AddParams struct {
	// URL is the URL of the site. If it doesn't contain a scheme
	// (like "http:" or "https:") it defaults to "https:".
	CameraID string `json:"cameraID"`
	FileData string `json:"fileData"`
}

// Add adds a new site to the list of monitored websites.
//
//encore:api public method=POST path=/motion_detection
func (s *Service) MotionDetection(context context.Context, params *AddParams) (*Motion, error) {
	// Prevent abuse by limiting the number of sites to 20.
	newImage, err := image_processing.Base64ToImage(params.FileData)
	if err != nil {
		return nil, &errs.Error{
			Code:    errs.InvalidArgument,
			Message: err.Error(),
		}
	}

	newImageGrey := image_processing.ImageToGrayscalePixels(newImage)
	oldImageGrey, ok := s.imageStore.GetLastFrame(params.CameraID)
	if !ok {
		s.imageStore.StoreFrame(params.CameraID, newImageGrey)
		return &Motion{Movement: false}, nil
	}

	result := image_processing.SubtractImages(newImageGrey, oldImageGrey)

	image_processing.StackBlur(&result, 2)

	image_processing.GlobalThresholding(result, 10)

	movement := image_processing.FindContour(result, 20)
	s.imageStore.StoreFrame(params.CameraID, newImageGrey)
	return &Motion{Movement: movement}, nil
}

// This is a service struct, learn more: https://encore.dev/docs/primitives/services-and-apis/service-structs
//
//encore:service
type Service struct {
	imageStore *image_processing.CameraStore
}

// initService is automatically called by Encore when the service starts up.
func initService() (*Service, error) {
	imageStore := image_processing.NewCameraStore()
	return &Service{imageStore: imageStore}, nil
}
