const progressBar = document.getElementById('progress-bar');
const loadingScreen = document.getElementById('loading-screen');
const iframe = document.getElementById('gameFrame');

function updateProgress(percent) {
    progressBar.style.width = percent + '%';
}

function fetchWithProgress(url) {
    return fetch(url).then(response => {
        const contentLength = response.headers.get('Content-Length');

        if (!contentLength) {
            throw new Error('Content-Length header is missing');
        }

        const total = parseInt(contentLength, 10);
        let loaded = 0;

        return new Response(
            new ReadableStream({
                start(controller) {
                    const reader = response.body.getReader();

                    function read() {
                        reader.read().then(({ done, value }) => {
                            if (done) {
                                controller.close();
                                return;
                            }

                            loaded += value.byteLength;
                            const percentComplete = Math.round((loaded / total) * 100);
                            updateProgress(percentComplete);

                            controller.enqueue(value);
                            read();
                        }).catch(error => {
                            console.error('Error reading stream:', error);
                            controller.error(error);
                        });
                    }

                    read();
                }
            })
        );
    });
}

fetchWithProgress('wasm/pongo.wasm').then(response => {
    return response.arrayBuffer();
}).then(bytes => {
    updateProgress(100);
    setTimeout(() => {
        loadingScreen.style.display = 'none';
        iframe.style.display = 'block';

        iframe.contentWindow.postMessage({ wasmBytes: bytes }, '*');
    }, 500);
}).catch(err => {
    console.error("Error loading WASM:", err);
});