<script lang="ts">
   import { onMount } from "svelte";
   import { pb } from "../utils/pb";
   import QRCode from "qrcode-generator";
   type Camera = {
      id: string;
      email: string;
      pictureResolution: string;
      detectionResolution: string;
      intervalBetweenPictures: string;
      sensitivity: string;
      [key: string]: any;
   };
   type Picture = {
      id: string;
      camera: string;
      image: string;
      created: string;
      [key: string]: any;
   };
   let cameras: Camera[] = $state([]);
   let selectedPictureIndex: number | null = $state(null);
   let selectedCameraPictures: Picture[] = $state([]);
   let fileToken: string | null = $state(null); // Store auth token
   let cameraName: string = $state("");
   let showQRCode: boolean = $state(false);
   let qrCodeData: string = $state("");
   // Fetch cameras and images on mount
// Fetch cameras and images on mount
onMount(async () => {
   if (!pb.authStore.isValid) {
      window.location.href="#/login";
      return
   }
   try {
      fileToken = await pb.files.getToken();
      const fetchedCameras = await pb.collection("cameras").getFullList();
      
      // Fetch all pictures in one request
      const allPictures = await pb.collection("pictures").getFullList(100, {
         sort: "-created",
      });
      
      // Map pictures to their respective cameras
      cameras = fetchedCameras.map(camera => {
         camera.pictures = allPictures.filter(pic => pic.camera === camera.id);
         return camera;
      });
      
   } catch (error) {
      console.error("Error fetching cameras:", error);
   }
});

   // Function to update camera in PocketBase
   async function updateCamera(cameraId: string, field: string, value: string) {
      const camera = cameras.find((c) => c.id === cameraId);
      if (!camera) return;

      // Update the local camera data
      camera[field] = value;

      try {
         // Use pb to update the camera in PocketBase
         await pb.collection("cameras").update(cameraId, {
            [field]: value,
         });
         console.log("Camera updated successfully");
      } catch (error) {
         console.error("Error updating camera:", error);
      }
   }
   async function createCamera() {
      if (!cameraName) {
         alert("Du musst einen KameraNamen eingeben");
         return;
      }
      let email = `${cameraName}@2gather.ch`;
      let password = crypto.randomUUID();

      try {
         // First create the camera
         await pb.collection("cameras").create({
            email: email,
            pictureResolution: "1920",
            detectionResolution: "144",
            intervalBetweenPictures: "1000",
            sensitivity: "500",
            password: password,
            passwordConfirm: password,
            emailVisibility: true,
         });

         // Get separate token for camera
         let res = await fetch(
            "http://127.0.0.1:8090/api/collections/cameras/auth-with-password",
            {
               method: "POST",
               headers: {
                  "Content-Type": "application/json",
               },
               body: JSON.stringify({
                  identity: email,
                  password: password,
               }),
            },
         );

         if (!res.ok) {
            alert("Something went wrong");
            return;
         }

         let data = await res.json();
         if (data.token) {
            showQRPopup(data.token);
         }
      } catch (error) {
         console.error("Error creating camera:", error);
         alert("Something went wrong");
      }
   }
   // Generate a secure URL for images
   function getSecureImageUrl(picture: Picture) {
      return `${pb.baseURL}/api/files/pictures/${picture.id}/${picture.image}?token=${fileToken}`;
   }

   // Open image viewer and set selected images
   function openImageViewer(cameraPictures: Picture[], index: number) {
      selectedCameraPictures = cameraPictures;
      selectedPictureIndex = index;
   }

   // Close image viewer
   function closeImageViewer() {
      selectedPictureIndex = null;
   }

   // Navigate to the previous image
   function prevImage() {
      if (selectedPictureIndex) {
         selectedPictureIndex--;
      }
   }

   // Navigate to the next image
   function nextImage() {
      if (
         selectedPictureIndex < selectedCameraPictures.length - 1
      ) {
         selectedPictureIndex++;
      }
   }

   // Close modal when Escape key is pressed
   function handleKeydown(event: { key: string }) {
      if (event.key === "Escape") {
         closeImageViewer();
      } else if (event.key === "ArrowLeft") {
         prevImage();
      } else if (event.key === "ArrowRight") {
         nextImage();
      }
   }

   // Add event listener for key presses
   $effect(() => {
      if (selectedPictureIndex !== null) {
         window.addEventListener("keydown", handleKeydown);
         return () => {
            window.removeEventListener("keydown", handleKeydown);
         };
      }
   });

   function showQRPopup(token: string) {
      // Create an object with the camera credentials
      const cameraData = {
         baseUrl: pb.baseURL,
         token: token,
      };

      // Convert to JSON string
      qrCodeData = JSON.stringify(cameraData);
      showQRCode = true;
   }

   // Function to generate QR code
   function generateQRCode(data: string) {
      const qr = QRCode(0, "L");
      qr.addData(data);
      qr.make();
      return qr.createDataURL(10);
   }

   // Function to close QR popup
   function closeQRPopup() {
      showQRCode = false;
      qrCodeData = "";
   }
</script>

<h1>Kameras</h1>
Neue Kamera<input bind:value={cameraName} /><button onclick={createCamera}
   >Erstellen</button
>
<div class="grid">
   {#each cameras as camera}
      <div class="camera-card">
         <h1>
            {camera.email}
         </h1>
         <div class="header">
            Auflösung:
            <input
               class="resolution"
               bind:value={camera.pictureResolution}
               oninput={() =>
                  updateCamera(
                     camera.id,
                     "pictureResolution",
                     camera.pictureResolution,
                  )}
            />
            Erkennungs-Auflösung:
            <input
               class="detection"
               bind:value={camera.detectionResolution}
               oninput={() =>
                  updateCamera(
                     camera.id,
                     "detectionResolution",
                     camera.detectionResolution,
                  )}
            />
            Interval(ms):
            <input
               class="detection"
               bind:value={camera.intervalBetweenPictures}
               oninput={() =>
                  updateCamera(
                     camera.id,
                     "intervalBetweenPictures",
                     camera.intervalBetweenPictures,
                  )}
            />
            Empfindlichkeit:
            <input
               class="detection"
               bind:value={camera.sensitivity}
               oninput={() =>
                  updateCamera(camera.id, "sensitivity", camera.sensitivity)}
            />
         </div>

         <div class="pictures">
            {#each camera.pictures as picture, i}
               <button
                  id={picture.id}
                  onclick={() => openImageViewer(camera.pictures, i)}
               >
                  <img
                     src={getSecureImageUrl(picture)}
                     alt={`Image from camera ${camera.id}`}
                  />
               </button>
            {/each}
         </div>
      </div>
   {/each}
</div>

<!-- Fullscreen Image Viewer -->
{#if selectedPictureIndex !== null}
   <div class="modal" onclick={closeImageViewer}>
      <div class="modal-content" onclick={(e) => e.stopPropagation()}>
         <button class="close-btn" onclick={closeImageViewer}>✖</button>
         {new Date(selectedCameraPictures[selectedPictureIndex].created)}
         <img
            src={getSecureImageUrl(
               selectedCameraPictures[selectedPictureIndex],
            )}
            alt="Fullscreen Image"
         />

         <!-- Navigation Buttons -->
         {#if selectedPictureIndex > 0}
            <button class="nav left" onclick={prevImage}>⬅</button>
         {/if}
         {#if selectedPictureIndex < selectedCameraPictures.length - 1}
            <button class="nav right" onclick={nextImage}>➡</button>
         {/if}
      </div>
   </div>
{/if}
{#if showQRCode}
   <div class="modal" onclick={closeQRPopup}>
      <div
         class="modal-content qr-content"
         onclick={(e) => e.stopPropagation()}
      >
         <button class="close-btn" onclick={closeQRPopup}>✖</button>
         <h2>Camera QR Code</h2>
         <div class="qr-container">
            <img src={generateQRCode(qrCodeData)} alt="Camera QR Code" />
         </div>
         <p class="qr-instructions">
            Scan this QR code with your camera app to configure it
         </p>
      </div>
   </div>
{/if}

<style>
   h1 {
      font-size: 24px;
      font-weight: bold;
      margin-bottom: 16px;
   }

   .grid {
      display: flex;
      flex-direction: column;
      gap: 16px;
      max-width: 1600px;
      margin: 0 auto;
   }

   .camera-card {
      width: 100%;
      border-radius: 8px;
      padding: 16px;
      box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
      border: 1px solid #ddd;
   }

   .header {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
      align-items: center;
      font-weight: bold;
      margin-bottom: 10px;
   }

   .pictures {
      display: flex;
      gap: 8px;
      overflow-x: auto;
      padding-top: 8px;
   }



   .pictures img {
      width: 80px;
      height: 80px;
      object-fit: cover;
      border-radius: 6px;
      border: 1px solid #ddd;
      cursor: pointer;
      transition: transform 0.2s;
   }

   .pictures img:hover {
      transform: scale(1.1);
   }

   /* Fullscreen Modal */
   .modal {
      position: fixed;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: rgba(0, 0, 0, 0.8);
      display: flex;
      align-items: center;
      justify-content: center;
   }

   .modal-content {
      position: relative;
      max-width: 80%;
      max-height: 90%;
   }

   .modal-content img {
      width: 100%;
      height: auto;
      max-height: 80vh;
      border-radius: 10px;
   }

   .close-btn {
      position: absolute;
      top: 10px;
      right: 15px;
      font-size: 24px;
      background: none;
      color: white;
      border: none;
      cursor: pointer;
   }

   .nav {
      position: absolute;
      top: 50%;
      font-size: 24px;
      background: rgba(255, 255, 255, 0.6);
      border: none;
      cursor: pointer;
      padding: 10px;
      border-radius: 50%;
      transform: translateY(-50%);
   }

   .nav.left {
      left: 10px;
   }

   .nav.right {
      right: 10px;
   }

   .nav:hover {
      background: white;
   }
</style>
