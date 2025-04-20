<script lang="ts">
  import QrScanner from "qr-scanner";

  // ✅ Get props dynamically
  let { onScan = (result:string) => console.log("Scanned:", result),isScanning=false } = $props();
  // ✅ Reactive state
  let scanResult = $state("");
  let scanning = $state(false);
  
  let videoElement: HTMLVideoElement;
  let scanner: QrScanner | null = null;

  export function startScanning() {
    if (!scanner) {
      scanner = new QrScanner(videoElement, (result) => {
        scanResult = result.data;
        onScan(result.data); // Call the passed function
        stopScanning(); // Stop after one successful scan
      }, {
        returnDetailedScanResult: true,
        highlightScanRegion: true,
      });
    }

    scanner.start()
      .then(() => scanning = true)
      .catch(err => console.error("QR Scanner Error:", err));
  }

  export function stopScanning() {
    if (scanner) {
      scanner.stop();
      scanning = false;
    }
  }

  // ✅ Cleanup
  $effect(() => {
    return () => {
      if (scanner) {
        scanner.destroy();
        scanner = null;
      }
    };
  });
</script>

<div>
  <video bind:this={videoElement} style="width: 100%; max-width: 1920px;">
    <track kind="captions" src="" srclang="en" label="No captions available">
  </video>
</div>
