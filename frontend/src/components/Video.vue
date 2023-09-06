<script setup>
    import { ref, onMounted } from 'vue';

    const stream = ref(null)
    const localStream = ref(null)
    const combinedStream = ref(null)
    const name = ref("")
    const overlay = ref("one")
    const color = ref("#220391")
    const textColor = ref("#f3f2f5")

    async function startStream(){
        try {
            stream.value = await navigator.mediaDevices.getUserMedia({ video: true, audio: true })
            localStream.value = await navigator.mediaDevices.getUserMedia({ video: true })
            const canvas = document.createElement('canvas');
            canvas.width = localStream.value.getVideoTracks()[0].getSettings().width;
            canvas.height = localStream.value.getVideoTracks()[0].getSettings().height;
            const ctx = canvas.getContext('2d');
            window.setInterval(()=>{
                ctx.drawImage(videoPlayer, 0, 0, canvas.width, canvas.height)
                if (overlay.value == "two"){
                    ctx.fillStyle = color.value
                    ctx.fillRect(0, canvas.height-50, canvas.width, 50)
                }
                ctx.font = '24px Arial'
                ctx.fillStyle = textColor.value
                const textWidth = ctx.measureText(name.value)
                ctx.fillText(name.value, (canvas.width-textWidth.width)/2, canvas.height-20)
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
    <div id="canvas" v-if="combinedStream !== null">
        <video :srcObject="combinedStream" autoplay></video>
    </div>
    <div id="button">
        <button  v-on:click="startStream">Start</button>
    </div>
    <div id="getname">
        <label for="name">Text on Overlay:</label>
        <input type="text" id="name" v-model="name">
    </div>
    <div id="options">
        <input type="radio" id="one" name="overlay" value="one" v-model="overlay">
        <label for="one">No Overlay</label>
        <input type="radio" id="two" name="overlay" value="two" v-model="overlay">
        <label for="two">Overlay</label>
    </div>
    <div id="color" v-if="overlay==='two'">
        <input type="color" id="color" v-model="color">
        <label for="color">Overlay Color</label>
        <input type="color" id="textColor" v-model="textColor">
        <label for="textColor">Text Color</label>
    </div>
</template>

<style>

div#video {
    display: none;
}

#videoPlayer {
    width: 500px;
    height: 300px;
}

</style>