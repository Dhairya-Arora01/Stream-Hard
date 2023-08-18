<script setup>
    import { ref, onMounted } from 'vue';

    const stream = ref(null)
    const localStream = ref(null)

    const streamConfig = {
        video: {
            width: { min: 178, ideal: 1280, max: 1920 },
            height: { min: 100, ideal: 720, max: 1080 },
            facingMode: 'user'
        },
        audio: true
    }

    const localStreamConfig = {
        video: {
            width: { min: 1024, ideal: 1280, max: 1920 },
            height: { min: 576, ideal: 720, max: 1080 },
            facingMode: 'user'
        }
    }

    // stun server
    const iceConfig = {'iceServers': [{'urls': 'stun:stun.l.google.com:19302'}]}


    onMounted(async() => {

        try {
            stream.value = await navigator.mediaDevices.getUserMedia(streamConfig)
            localStream.value = await navigator.mediaDevices.getUserMedia(localStreamConfig)
        } catch (error) {
            console.error("Webcam not working", error)
        }

    })

    async function iceServer() {

        try {
            const peerConnection = new RTCPeerConnection(iceConfig)
            stream.value.getTracks().forEach(track => {
                peerConnection.addTrack(track, stream.value)
            })

            const offer = await peerConnection.createOffer()
            await peerConnection.setLocalDescription(offer)
            // console.log('local Description',peerConnection.localDescription)
            sendOfferToServer(peerConnection)
        } catch (error) {
            console.error("Ice fault",error)
        }

    }

    async function sendOfferToServer(peerConnection) {

        const result = await fetch("http://localhost:8000/feed", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                'sdp': peerConnection.localDescription,
            })
        })

        const msg = await result.json()
        console.log("message", msg)

    }

</script>

<template>
    <div id="video">
        <video :srcObject="localStream" id="videoPlayer" autoplay></video>
    </div>
    <div id="button">
        <button v-on:click="iceServer">Start</button>
    </div>
</template>

<style>

div#video {

    width: 600px;
    height: 400px;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    background-color: blue;

}

#videoPlayer {
    width: 500px;
    height: 300px;
}

</style>