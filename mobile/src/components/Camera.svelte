<script lang="ts">
	import { pb } from '$lib';
	import { onDestroy, onMount } from 'svelte';

	let config = $state({
        detectionResolution: 120,
        imageResolution: 1920,
        interval: 1000});
	let video_source: HTMLVideoElement = $state()!;
	let movement = $state(false);
	let streamRunning = false;
	let canvas: HTMLCanvasElement;

	onMount(() => {
		
		getAuthData();
		canvas = document.createElement('canvas');
		getConfig();
	});
	async function getConfig() {
		if (!pb.authStore?.record?.id) return;
		const data= await pb
			.collection('cameras')
			.getOne(pb.authStore.record.id, { fields: 'detectionResolution,imageResolution,interval' }) as any;
		config = data;
		pb.collection('cameras').subscribe(
			pb.authStore.record.id,
			function (e) {
				console.log('new config: ', e.record);
				config = e.record as any;
			},
			{ fields: 'detectionResolution,imageResolution,interval' }
		);
	}

	onDestroy(() => (streamRunning = false));
	async function access_webcam() {
		try {
			const stream = await navigator.mediaDevices.getUserMedia({
				video: true
			});

			video_source.srcObject = stream;
			video_source.play();
			streamRunning = true;
			video_source.addEventListener('ended', () => {
				streamRunning = false;
			});
			video_source.addEventListener('error', () => {
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
		takePicture(config.detectionResolution);
		setTimeout(takePictureInterval, config.interval);
	}

	async function takePicture(imageWidth: number) {
		const context = canvas.getContext('2d')!;
		if (!video_source.videoWidth && video_source.videoHeight) {
			return;
		}
		let ratio = video_source.videoWidth / video_source.videoHeight;
		let imageHeight = Math.round(imageWidth / ratio);
		canvas.width = imageWidth;
		canvas.height = imageHeight;
		context.drawImage(video_source, 0, 0, imageWidth, imageHeight);
		let data = canvas.toDataURL('image/jpg');

		// Send the image to the motion detection API
		const res = await pb.send('/api/motion_detection', {
			method: 'POST',
			body: JSON.stringify({
				imageData: data.split(';base64,')[1], // Optional - Base64 part for motion detection (if needed)
				cameraID: pb.authStore.record?.id
			})
		});

		if (res.movement) {
			await uploadHighResolutionImage(config.imageResolution);
		}
	}
	async function uploadHighResolutionImage(width: number) {
		// Create a new canvas for high-resolution capture
		const highResCanvas = document.createElement('canvas');
		const highResContext = highResCanvas.getContext('2d')!;

		// Calculate the new height based on the aspect ratio
		const aspectRatio = video_source.videoWidth / video_source.videoHeight;
		const newHeight = Math.round(width / aspectRatio);

		// Set the canvas size to the new dimensions
		highResCanvas.width = width;
		highResCanvas.height = newHeight;

		// Draw the video on the high-res canvas
		highResContext.drawImage(video_source, 0, 0, width, newHeight);

		// Create a blob from the full-size (high-res) canvas
		const highResBlob: any = await new Promise((resolve) =>
			highResCanvas.toBlob(resolve, 'image/jpeg')
		);

		// Create a FormData object for high-resolution upload
		const highResFormData = new FormData();
		highResFormData.append('file', highResBlob, `${Date.now()}.jpg`); // 'file' is the expected key for file uploads

		// Send the high-resolution file data to PocketBase to create a new record
		await pb.collection('pictures').create({
			image: highResFormData.get('file'), // Pass the file directly from FormData
			camera: pb.authStore.record?.id
		});
	}
	async function getAuthData() {
console.log('getAuthData')
	}
	access_webcam();
</script>

	<!-- svelte-ignore a11y_media_has_caption -->
	<video bind:this={video_source} class=" h-3/4 w-full"></video>
	{#if movement}
		<div>
			<h6>Movement Detected</h6>
		</div>
	{/if}


