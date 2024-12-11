<script lang="ts">
    import { onDestroy, onMount } from "svelte";

    let imageWidth = 120;
    let interval=200;
    let video_source: HTMLVideoElement = $state()!;
    let movement = $state(false);
    let streamRunning = false;
    let canvas: HTMLCanvasElement;
    onMount(() => (canvas = document.createElement("canvas")));
    onDestroy(() => streamRunning=false);
    async function access_webcam() {
        try {
            const stream = await navigator.mediaDevices.getUserMedia({
                video: true,
            });

            video_source.srcObject = stream;
            video_source.play();
            streamRunning = true;
            video_source.addEventListener("ended", () => {
                streamRunning = false;
            });
            video_source.addEventListener("error", () => {
                streamRunning = false;
            });
            takePictureInterval();
        } catch (error) {
            streamRunning = false;
            console.error(error);
        }
    }

    function takePictureInterval() {
        if (!streamRunning) return;
        takePicture();
        setTimeout(takePictureInterval, interval);
    }

    async function takePicture() {
        const context = canvas.getContext("2d")!;
        if (video_source.videoWidth && video_source.videoHeight) {
            let ratio = video_source.videoWidth / video_source.videoHeight;
            let imageHeight = Math.round(imageWidth / ratio);
            canvas.width = imageWidth;
            canvas.height = imageHeight;
            context.drawImage(video_source, 0, 0, imageWidth, imageHeight);
            let data = canvas.toDataURL("image/jpg");
            data = data.split(";base64,")[1];
            let response = await fetch("/motion_detection", {
                method: "POST",
                body: JSON.stringify({
                    fileData: data,
                    cameraID: "1234",
                }),
            });
            const responseData = await response.json();
            movement = responseData.movement;

        }
    }

    access_webcam();

</script>

<div class="h-full w-full relative">
    <!-- svelte-ignore a11y_media_has_caption -->
    <video bind:this={video_source} class=" h-full w-full"></video>
    {#if movement}
        <div>
            <h6>Movement Detected</h6>
        </div>
    {/if}
</div>

<style lang="postcss">
    video {
        -webkit-transform: scaleX(-1);
        transform: scaleX(-1);
    }
</style>
