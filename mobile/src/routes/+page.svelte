<script lang="ts">
	import { onMount } from 'svelte';
	import QrScanner from '../components/QRScanner.svelte';
	import Camera from '../components/Camera.svelte';
	import { pb } from '$lib';
	let qrScanner: QrScanner;
		onMount(async () => {
		try {
			let b=await pb.collection('cameras').authRefresh();
		} catch (error) {
			console.error('Error refreshing token:', error);
		}
	
	
		qrScanner.startScanning();
	});

	function onQRScan(result: string) {
		qrScanner.stopScanning();
		const data=JSON.parse(result);
		if(!data.token) return;
		pb.authStore.save(data.token);
	
		qrScanner.stopScanning();
		window.location.reload();
	}
	function logout() {
		pb.authStore.clear();
		window.location.reload();
	}
</script>
<div class="bg-gray-900 text-white flex flex-col items-center justify-center h-screen">
{#if !pb.authStore.isValid}
<QrScanner bind:this={qrScanner} onScan={onQRScan} />
{/if}
{#if pb.authStore.isValid}
<button onclick={logout} class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Logout</button>
<Camera />
{/if}
</div>


