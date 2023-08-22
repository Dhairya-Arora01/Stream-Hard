<script setup>
    import { ref, onMounted } from 'vue';

    const stream = ref(null)

    async function startStream(){
        try {
            stream.value = await navigator.mediaDevices.getUserMedia({ video:true })
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
        stream.value.getTracks().forEach(track => peer.addTrack(track, stream.value));

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
        <video :srcObject="stream" id="videoPlayer" autoplay></video>
    </div>
    <div id="button">
        <button  v-on:click="startStream">Start</button>
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