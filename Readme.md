# Stream-Hard
A live streaming platform to broadcast live video from your browser to various social media platforms.<hr>

## Features 💎
- ⚡ Live stream your webcam video from your browser to social media platforms.
- 🫂 User friendly web-interface.
- 🔐 Authentication for clients.
- ✨ Add overlays to your stream.
- 🙌 Supports all platforms that allows streaming through an rtmp link.
- 🚫 We do not store your rtmp links and keys.

## Getting Started 🚀
**Prerequesite** : You need to have docker and docker-compose installed on your machine.
1. Clone the repository on your machine.
~~~
git clone https://github.com/Dhairya-Arora01/Stream-Hard.git
~~~
2. cd into the cloned repo.
3. Start the application components using
~~~
docker compose up
~~~
5. Now you can access the application at localhost:8080 in your browser.
6. Create Account, enter your rtmp link and start streaming.
   ![Stream-Hard demo](https://i.imgur.com/hK6r6Ih.gif)
7. [Here](https://i.imgur.com/Ny2bep2.jpeg) is the screenshot of the output stream on twitch.

## Architecture 🏛
- **Frontend** : Vue
- **Backend** : Go
- **Database** : Postgres
- **Other technologies used**:
  - pion/webrtc for WebRTC
  - ffmpeg for sendting the stream to rtmp servers.
    
![Architecture](https://i.imgur.com/miWGRow.jpeg)

## Under the hood ⚙️
After authentication, the user enters the rtmp link and hits go. Then,
1. This sends request to the backend server.
2. This opens a websocket between the client and the backend.
3. Both of them asynchronously obtain ICE candidates from the stun server.
4. Through the opened websocket they exchange ICE candidates, offers, answers and simultaneously set local and remote session descriptions.
5. WebRTC connection is established !
6. Now the webcam stream in the form RTP packets is sent directly to the backend.
7. ffmpeg forwards the stream to the provided RTMP server.

## At Scale..📈
Currently we are also developing the Kubernetes repository for our Stream-Hard application that contains all k8s specific configurations. You can find the repository at https://github.com/Dhairya-Arora01/Stream-Hard-config . Your contributions are welcome on the configuration repository as well.

## Support and Issues 💙
If you encounter any issues or have questions, feel free to open a new issue.

## Contributing 💝
We welcome contributions from the community. Feel free to submit pull requests or open issues.

## Acknowledgments 🙏
- [Pion webrtc](https://github.com/pion/webrtc)
- [ffmpeg](https://github.com/FFmpeg/FFmpeg)
- [Sean Dubois](https://github.com/Sean-Der) - Creator of Pion WebRTC and an inspiration.

## Tasks checklist 📝
- [x] Vue basic
- [x] Go basic
- [x] WebRTC connection between browser and server
- [x] Send webcam video to server using webRTC
- [x] Streams transfer to rtmp using ffmpeg
- [x] Cosmetic Frontend
- [x] Input for rtmp link
- [x] Adding overlays
- [x] User Authentication
- [x] Concurrency in backend / error resillience
- [x] Implementing Docker
- [ ] Local k8s cluster setup - ....In progress....
- [ ] IAC to setup the cluster in the cloud
- [ ] Unit and e2e testing
- [ ] CI using Github Actions
- [ ] ArgoCD for GitOps
- [ ] Monitoring and Logging
