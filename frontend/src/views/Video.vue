<script setup>
    import { ref, onMounted } from 'vue';
import router from '../router';

    const socket = ref(null)
    const peer = ref(null)
    const audioStream = ref(null)
    const localStream = ref(null)
    const combinedStream = ref(null)
    const rtmpLink = ref("")
    const name = ref("")
    const overlay = ref("one")
    const color = ref("#220391")
    const textColor = ref("#f3f2f5")
    const active = ref(false)

    async function startStream(){
        const token = localStorage.getItem("token")

        if (token == null) {
            router.replace('/login')
        }

        try {     
            audioStream.value = await navigator.mediaDevices.getUserMedia({ audio: true })
            localStream.value = await navigator.mediaDevices.getUserMedia({ video: true })

            const canvas = document.createElement('canvas');
            canvas.width = localStream.value.getVideoTracks()[0].getSettings().width
            canvas.height = localStream.value.getVideoTracks()[0].getSettings().height
            const ctx = canvas.getContext('2d')

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

            const canvasStream = canvas.captureStream()
            combinedStream.value = new MediaStream([...audioStream.value.getTracks(), ...canvasStream.getTracks()])

        } catch (error) {
            console.error("Webcam not working", error)
        }
        
        socket.value = new WebSocket("ws://localhost:8000/ws?bearer=" + token)
        socket.value.onerror = e => {
            console.error("Error opening websocket returning you back!")
        }

        peer.value = new RTCPeerConnection({
            iceServers: [
                {
                    urls: ["stun:stun1.l.google.com:19302"]
                },
            ],
        })
        combinedStream.value.getTracks().forEach(track => peer.value.addTrack(track, combinedStream.value));

        socket.value.onmessage = e =>{
            let msg = JSON.parse(e.data)
            if (!msg){
                return console.log("failed to parse msg")
            }

            if (msg.candidate) {
                peer.value.addIceCandidate(msg)
            } else if (msg.type){
                peer.value.setRemoteDescription(msg)
            } else if (msg.RTMPError){
                console.log(msg.RTMPError)
                endStream()
            }
        }

        socket.value.onopen = ()=>{
            socket.value.send(JSON.stringify({rtmp: rtmpLink.value}))
            active.value = true
            peer.value.createOffer().then(offer => {
                peer.value.setLocalDescription(offer)
                socket.value.send(JSON.stringify(offer))
            })
        }

        socket.value.onclose = ()=> {
            endStream()
        }
        
    }

    async function endStream() {
        peer.value.close()
        socket.value.close()
        active.value = false
    }

</script>

<template>
    <div id="video">
        <video :srcObject="localStream" id="videoPlayer" autoplay></video>
    </div>
    <div id="canvas" v-if="combinedStream !== null">
        <video :srcObject="combinedStream" autoplay muted></video>
    </div>
    <div id="button">
        <button  v-on:click="startStream" v-if="!active">Start</button>
        <button v-on:click="endStream" v-if="active">Stop</button>
    </div>
    <div id="rtmp-link">
        <label for="rtmp-link">rtmp:</label>
        <input type="text" id="rtmp-link" v-model="rtmpLink">
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