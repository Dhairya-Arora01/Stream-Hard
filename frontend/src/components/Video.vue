<script setup>
    import { ref, onMounted } from 'vue';

    const stream = ref(null)
    const localStream = ref(null)
    const combinedStream = ref(null)
    const name = ref(null)

    async function startStream(){
        try {
            stream.value = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
            localStream.value = await navigator.mediaDevices.getUserMedia({ video: true })
            const canvas = document.createElement('canvas');
            canvas.width = localStream.value.getVideoTracks()[0].getSettings().width; // Adjust based on your video size
            canvas.height = localStream.value.getVideoTracks()[0].getSettings().height; // Adjust based on your video size
            const ctx = canvas.getContext('2d');
            window.setInterval(()=>{
                ctx.drawImage(videoPlayer, 0, 0, canvas.width, canvas.height)
                ctx.font = '24px Arial'
                ctx.fillStyle = 'black'
                ctx.fillText(name.value, 10, 30)
            })
            combinedStream.value = canvas.captureStream();
        } catch (error) {
            console.error("Webcam not working", error)
        }

        const socket = new WebSocket("ws://localhost:8000/ws")
        const peer = new RTCPeerConnection({
            iceServers: [
                {
                    urls: ["stun:stun1.l.google.com:19302"]
                },
            ],
        })
        combinedStream.value.getTracks().forEach(track => peer.addTrack(track, combinedStream.value));

        socket.onmessage = e =>{
            let msg = JSON.parse(e.data)
            if (!msg){
                return console.log("failed to parse msg")
            }

            if (msg.candidate) {
                peer.addIceCandidate(msg)
            } else if (msg.type){
                peer.setRemoteDescription(msg)
            }
        }

        socket.onopen = ()=>{
            peer.createOffer().then(offer => {
                peer.setLocalDescription(offer)
                socket.send(JSON.stringify(offer))
            })
        }
        
    }

</script>

<template>
    <div id="video">
        <video :srcObject="localStream" id="videoPlayer" autoplay></video>
    </div>
    <div id="button">
        <button  v-on:click="startStream">Start</button>
    </div>
    <div id="getname">
        <input type="text" v-model="name">
    </div>
    <div>
        <p>canvas</p>
        <video :srcObject="combinedStream" autoplay></video>
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