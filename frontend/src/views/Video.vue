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
    const username = ref("")

    onMounted(()=>{
        const token = localStorage.getItem("token")

        if (token == null) {
            router.replace('/login')
        }

        username.value = localStorage.getItem("name")
    })

    async function startStream(){
        

        const token = localStorage.getItem("token")
        try {     
            audioStream.value = await navigator.mediaDevices.getUserMedia({ audio: true })
            localStream.value = await navigator.mediaDevices.getUserMedia({ video: true })

            const canvas = document.createElement('canvas');
            canvas.width = localStream.value.getVideoTracks()[0].getSettings().width
            canvas.height = localStream.value.getVideoTracks()[0].getSettings().height
            const ctx = canvas.getContext('2d')

            const videoPlayer = document.getElementById("videoPlayer")

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
            store.commit("funcErr", "Webcam not working!")
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
                store.commit("funcErr", msg.RTMPError)
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
            store.commit("funcErr", "Stream Ended!")
            endStream()
        }
        
    }

    async function endStream() {
        peer.value.close()
        socket.value.close()
        active.value = false
    }

    function logout() {
        localStorage.removeItem('token')
        router.replace('/')
    }

</script>

<template>
    <div id="main">
        <div id="greetings">
            <h3>Hello {{ username }},</h3>
            <button id="logout" v-on:click="logout">Logout</button>
        </div>
        <div id="screen">
            <div id="left">
                <div id="video">
                    <video :srcObject="localStream" id="videoPlayer" autoplay></video>
                </div>
                <div id="photo" v-if="combinedStream == null">
                    <img src="../assets/user.png" alt="user">
                </div>
                <div id="photo" v-if="combinedStream !== null">
                    <video :srcObject="combinedStream" autoplay muted></video>
                </div>
            </div>
            <div id="right">
                <div id="button">
                    <button id="start" v-on:click="startStream" v-if="!active">Start</button>
                    <button id="stop" v-on:click="endStream" v-if="active">Stop</button>
                </div>
                <div id="formelement">
                    <label for="rtmp-link">RTMP link</label>
                    <input type="text" id="rtmp-link" v-model="rtmpLink" required>
                </div>
                <div id="formelement">
                    <label for="name">Text on Overlay</label>
                    <input type="text" id="name" v-model="name">
                </div>
                <div id="options">
                    <input type="radio" id="one" name="overlay" value="one" v-model="overlay">
                    <label for="one">No Overlay</label>
                    <input type="radio" id="two" name="overlay" value="two" v-model="overlay">
                    <label for="two">Overlay</label>
                </div>
                <div id="options" v-if="overlay==='two'">
                    <input type="color" id="color" v-model="color">
                    <label for="color">Overlay Color</label>
                    <input type="color" id="textColor" v-model="textColor">
                    <label for="textColor">Text Color</label>
                </div>
            </div>
        </div>
        
    </div>
</template>

<style>

div#main {
    width: 80%;
    height: 100%;
    color: white;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    align-items: center;
    font-family: 'Poppins', sans-serif;
}

div#greetings {
    width: 100%;
    display: flex;
    flex-direction: row;
    justify-content: space-between;
}

div#greetings > h3 {
    font-size: 1.4rem;
}

button#logout {
    width: 6%;
    height: 100%;
    background-color: rgb(99, 106, 116);
    border: none;
    color: white;
}

button#logout:hover {
    background-color: rgb(126, 133, 144);
}

div#screen {
    display: flex;
    flex-direction: row;
    justify-content: space-around;
    align-items: center;
}

div#left {
    display: flex;
    justify-content: left;
    align-items: center;
}

div#right {
    height: 100%;
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    align-items:center;
}

div#video {
    display: none;
}

div#photo{
    height: 88%;
    width: 78%;
}

div#photo>img{
    width: 100%;
    height: 100%;
}

div#photo>video{
    width: 100%;
    height: 100%;
}

div#button {
    width: 6.2rem;
    height: 3rem;
}

div#button > button {
    width: 100%;
    height: 100%;
    font-size: 1.2rem;
    border: none;
    border-radius: 2px;
}

button#start {
    background-color: rgb(6, 242, 6);
}

button#start:hover {
    background-color: rgb(29, 180, 29);
}

button#stop {
    background-color: rgb(198, 29, 29);
}

button#stop:hover {
    background-color: rgb(149, 25, 25);
}

div#formelement {
    display: flex;
    flex-direction: column;
    margin: 2%;
    width: 100%;
}

div#formelement > input {
    height: 2.8em;
    font-size: 0.9em;
    color: white;
    border-radius: 0.2em;
    border: none;
    padding-left: 0.4em;
    background-color: rgb(99, 106, 116);
}

div#options {
    display: flex;
    flex-direction: row;
}

div#options {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    align-items: center;
    height: 1.5rem;
    color: white;
    margin: 1rem;
    width: 50%;
}

div#options input {
    margin: 2rem;
    height: 1.1rem;
}

</style>