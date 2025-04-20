package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"bee-aware/image_processing"
	_ "bee-aware/migrations"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
	app := pocketbase.New()

	imageStore := image_processing.NewCameraStore()
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/frontend/{path...}", apis.Static(os.DirFS("./frontend/dist/"), false))
		se.Router.GET("/mobile/{path...}", apis.Static(os.DirFS("./mobile/build/"), false))

		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))
		// register "POST /api/myapp/settings" route (allowed only for authenticated users)
		se.Router.POST("/api/motion_detection", func(e *core.RequestEvent) error {
			data := struct {
				ImageData string `json:"imageData" form:"imageData"`
				CameraID  string `json:"cameraID" form:"cameraID"`
			}{}
			if err := e.BindBody(&data); err != nil {
				return e.BadRequestError("Failed to read request data", err)
			}

			camera, err := app.FindRecordById("cameras", data.CameraID)
			if err != nil {
				return e.BadRequestError("Failed to fetch camera record", err)
			}

			newImage, err := image_processing.Base64ToImage(data.ImageData)
			if err != nil {
				return e.BadRequestError("Failed to parse image data", err)
			}

			newImageGrey := image_processing.ImageToGrayscalePixels(newImage)
			oldImageGrey, ok := imageStore.GetLastFrame(data.CameraID)
			if !ok {
				imageStore.StoreFrame(data.CameraID, newImageGrey)
				return e.JSON(http.StatusOK, map[string]bool{"movement": false})
			}

			result := image_processing.SubtractImages(newImageGrey, oldImageGrey)

			image_processing.StackBlur(&result, 2)

			image_processing.GlobalThresholding(result, 10)

			movement := image_processing.FindContour(result, camera.GetInt("sensitivity"))
			imageStore.StoreFrame(data.CameraID, newImageGrey)

			return e.JSON(http.StatusOK, map[string]bool{"movement": movement})
		}).Bind(apis.RequireAuth())

		return se.Next()
	})
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
