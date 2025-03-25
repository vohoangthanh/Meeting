import CloudflareCalls from './CloudflareCalls.js';

// Lấy current URL và cấu hình endpoints
const currentUrl = window.location.origin;
const isLocalhost = currentUrl.includes('localhost') || currentUrl.includes('127.0.0.1');

// Cấu hình endpoints dựa trên môi trường
const config = {
    backendUrl: isLocalhost ? 'http://localhost:60000' : currentUrl,
    websocketUrl: isLocalhost
        ? 'ws://localhost:50000/ws'
        : currentUrl.replace('http', 'ws').replace('https', 'wss') + '/ws',
    backendUrl2: isLocalhost ? 'http://localhost:60000' : currentUrl,
    websocketUrl2: isLocalhost
        ? 'ws://localhost:60000/ws'
        : currentUrl.replace('http', 'ws').replace('https', 'wss') + '/ws',
};
const baseAPI = isLocalhost ? 'http://localhost:60000' : currentUrl;
const calls = new CloudflareCalls(config);

// Tương tự cho screenShareCalls
const screenShareConfig = {...config};

let currentRoom = null;
let screenShareCalls = null;

// Add at the top with other globals
let isMaskEnabled = false;
let isBlurEnabled = false;
let currentMask = 'basic/mask1.png';
let masksList = [];
let faceMaskFilter = null;
let backgroundBlur = null;
let processedStream = null;


// DOM Elements
const videoGrid = document.getElementById('videoGrid');
const controls = {
    toggleMic: document.getElementById('toggleMicBtn'),
    toggleVideo: document.getElementById('toggleVideoBtn'),
    shareScreen: document.getElementById('shareScreenBtn'),
    toggleMask: document.getElementById('toggleMaskBtn'),
    wave: document.getElementById('waveBtn'),
    leave: document.getElementById('leaveBtn')
};
const participantsList = document.getElementById('participantsList');
const notificationsContainer = document.getElementById('notificationsContainer');

// Get stored data from localStorage
const username = localStorage.getItem('username');
const roomId = localStorage.getItem('roomId');

// Get token and initialize calls
async function ensureInitialized() {
    if (!calls.token) {
        try {
            calls.inauthtoken(username);
            // const response = await fetch(`${baseAPI}/auth/token`, {
            //     method: 'POST',
            //     headers: { 'Content-Type': 'application/json' },
            //     body: JSON.stringify({ username })
            // });

            // const { token } = await response.json();
            // calls.setToken(token);
            // showNotification('Successfully initialized');
            // return true;
        } catch (err) {
            console.error('Error getting token:', err);
            showNotification('Failed to initialize', 'error');
            return false;
        }
    }
    return true;
}

async function setupLocalVideo() {
    try {
        const stream = await navigator.mediaDevices.getUserMedia({
            video: true,
            audio: true
        });

        // Tìm hoặc tạo container cho local video
        let container = document.querySelector('.local-video');
        if (!container) {
            container = document.createElement('div');
            container.className = 'video-container local-video';

            // Tạo video element cho preview local
            const video = document.createElement('video');
            video.autoplay = true;
            video.playsInline = true;
            video.muted = true; // Đảm bảo video local luôn bị mute

            const nameLabel = document.createElement('div');
            nameLabel.className = 'participant-name';
            nameLabel.textContent = username || 'You';

            container.appendChild(video);
            container.appendChild(nameLabel);
            videoGrid.appendChild(container);
        }

        // Set stream nhưng đảm bảo audio luôn bị mute cho local preview
        const video = container.querySelector('video');
        video.srcObject = stream;
        video.volume = 0; // Thêm dòng này để đảm bảo volume = 0
        video.muted = true; // Thêm dòng này để doubly sure

        // Lưu stream cho WebRTC
        calls.localStream = stream;
        console.log('Local video setup complete');
    } catch (err) {
        console.error('Error accessing media devices:', err);
        showNotification('Failed to access camera/microphone', 'error');
    }
}

function getParticipantDisplayName(participant) {
    return participant.name || `User-${participant.userId.slice(0, 6)}`;
}

async function joinRoom() {
    try {
        // 1. Join room và lấy session
        await calls.joinRoom(roomId, { name: username });
        currentRoom = roomId;
        showNotification(`Joined room: ${roomId}`);

        // 2. Setup handlers trước khi pull tracks
        setupCallbacks();

        // 3. Lấy danh sách người tham gia
        const participants = await calls.listParticipants();
        console.log('Current participants:', participants);

        // 4. Set up remote streams for existing participants
        for (const participant of participants) {
            // Skip if it's our own session
            if (participant.sessionId === calls.sessionId) continue;

            console.log('Processing participant:', participant);

            // Create container for this participant if not exists
            const containerId = `participant-${participant.sessionId}`;
            if (!document.getElementById(containerId)) {
                const container = document.createElement('div');
                container.id = containerId;
                container.className = 'video-container';

                const video = document.createElement('video');
                video.autoplay = true;
                video.playsInline = true;

                const name = document.createElement('div');
                name.className = 'participant-name';
                name.textContent = participant.name || `participant-${participant.sessionId}`;

                container.appendChild(video);
                container.appendChild(name);
                videoGrid.appendChild(container);

                // Set up MediaStream for this participant
                video.srcObject = new MediaStream();
            }

            // Pull each track from the participant
            for (const trackName of participant.publishedTracks) {
                console.log(`Pulling track ${trackName} from session ${participant.sessionId}`);
                await calls._pullTracks(participant.sessionId, trackName);
            }
        }

        // 5. Start monitoring stats sau khi mọi thứ đã setup
        calls.startStatsMonitoring(1000);
    } catch (err) {
        console.error('Error joining room:', err);
        showNotification('Failed to join room: ' + err.message, 'error');
    }
}

async function setupScreenShare() {

    // Tạo một instance CloudflareCalls mới cho screen share
    screenShareCalls = new CloudflareCalls(screenShareConfig);

    // Khởi tạo token cho screen share
    const response = await fetch(`${baseAPI}/auth/token`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: username + '-screen' })
    });

    const { token } = await response.json();
    screenShareCalls.setToken(token);

    // Lấy screen share stream
    const screenStream = await navigator.mediaDevices.getDisplayMedia({
        video: true,
        audio: false
    });

    // Lưu stream và join room
    screenShareCalls.localStream = screenStream;
    await screenShareCalls.joinRoom(roomId, { name: username + "'s Screen" });

    // THÊM: Publish tracks sau khi join room
    await screenShareCalls.publishTracks();

    // Setup callback cho screen share để nhận remote tracks
    screenShareCalls.onRemoteTrack((track) => {
        console.log('Screen share track received:', track);
        const containerId = `participant-${track.sessionId}`;
        let container = document.getElementById(containerId);

        if (!container) {
            container = document.createElement('div');
            container.id = containerId;
            container.className = 'video-container';

            const video = document.createElement('video');
            video.autoplay = true;
            video.playsInline = true;

            const name = document.createElement('div');
            name.className = 'participant-name';
            name.textContent = username + "'s Screen";

            container.appendChild(video);
            container.appendChild(name);
            videoGrid.appendChild(container);

            video.srcObject = new MediaStream();
        }

        const video = container.querySelector('video');
        video.srcObject.addTrack(track);
    });

    // Lắng nghe sự kiện kết thúc screen share
    const screenTrack = screenStream.getVideoTracks()[0];
    screenTrack.onended = async () => {
        await stopScreenShare();
    };

    // Cập nhật icon
    controls.shareScreen.querySelector('.material-icons').textContent = 'stop_screen_share';
    showNotification('Screen sharing started');


}

async function stopScreenShare() {
    if (screenShareCalls) {
        // Dọn dẹp stream
        if (screenShareCalls.localStream) {
            screenShareCalls.localStream.getTracks().forEach(track => track.stop());
        }

        // Rời phòng và đóng kết nối
        await screenShareCalls.leaveRoom();
        screenShareCalls = null;

        // Cập nhật UI
        controls.shareScreen.querySelector('.material-icons').textContent = 'screen_share';
        showNotification('Screen sharing ended');
    }
}

function setupCallbacks() {
    calls.onRemoteTrack((track) => {
        console.log('Remote track received:', track);
        const containerId = `participant-${track.sessionId}`;
        let container = document.getElementById(containerId);

        if (!container) {
            container = document.createElement('div');
            container.id = containerId;
            container.className = 'video-container';

            const video = document.createElement('video');
            video.autoplay = true;
            video.playsInline = true;

            const name = document.createElement('div');
            name.className = 'participant-name';
            name.textContent = 'Participant ' + track.sessionId;

            container.appendChild(video);
            container.appendChild(name);
            videoGrid.appendChild(container);

            
        }

        const video = container.querySelector('video');
        if (!video.srcObject) {
            video.srcObject = new MediaStream();
        }
        video.srcObject.addTrack(track);
    });

    calls.onRemoteTrackUnpublished((sessionId, trackName) => {
        console.log('Remote track unpublished:', { sessionId, trackName });
        const containerId = `participant-${sessionId}`;
        const container = document.getElementById(containerId);

        if (container) {
            const video = container.querySelector('video');
            if (video && video.srcObject) {
                // Tìm và xóa track khỏi MediaStream
                const stream = video.srcObject;
                const tracks = stream.getTracks();
                tracks.forEach(track => {
                    if (track.id === trackName) {
                        stream.removeTrack(track);
                        track.stop();
                    }
                });

                // Nếu không còn track nào, xóa container
                if (stream.getTracks().length === 0) {
                    container.remove();
                }
            }
        }
    });

    // Thêm callback xử lý khi track status thay đổi
    calls.onTrackStatusChanged(async ({ sessionId, trackName, status }) => {
        console.log('Track status changed:', { sessionId, trackName, status });

        // Tìm container của participant
        const containerId = `participant-${sessionId}`;
        const container = document.getElementById(containerId);

        if (container) {
            const video = container.querySelector('video');
            if (video && video.srcObject) {
                // Tìm track cần cập nhật
                const mediaStream = video.srcObject;
                const tracks = mediaStream.getTracks();

                // Nếu track bị disabled, pull lại track mới
                if (status === 'disabled') {
                    // Xóa track cũ
                    tracks.forEach(track => {
                        if (track.id === trackName) {
                            mediaStream.removeTrack(track);
                            track.stop();
                        }
                    });

                    // Pull track mới
                    try {
                        await calls._pullTracks(sessionId, trackName);
                        console.log(`Re-pulled track ${trackName} for session ${sessionId}`);
                    } catch (error) {
                        console.error('Error re-pulling track:', error);
                    }
                }
            }
        }
    });

    calls.onParticipantLeft((participant) => {
        const container = document.getElementById(`participant-${participant.sessionId}`);
        if (container) {
            container.remove();
        }
        showNotification(`${participant.name || 'A participant'} left the room`);

    });

    // Sửa lại handler cho data messages
    calls.onDataMessage(async (data) => {
        console.log('Received data message:', data);
        showNotification(`👋 ${data.data.fromName} vẫy tay chào!`);
    });

    // Sửa lại handler cho nút vẫy tay
    controls.wave.onclick = async () => {
        try {
            // Gửi tin nhắn vẫy tay với đầy đủ thông tin
            await calls.sendDataToAll({
                type: 'wave',
                fromSession: calls.sessionId,
                fromName: username,
                timestamp: Date.now() // Thêm timestamp để tránh trùng lặp
            });

            // Hiển thị thông báo local
            showNotification('Bạn đã vẫy tay chào mọi người! 👋');

        } catch (error) {
            console.error('Error sending wave:', error);
            showNotification('Không thể gửi vẫy tay', 'error');
        }
    };

    // Control buttons
    controls.toggleMic.onclick = () => {
        const audioTrack = calls.localStream?.getAudioTracks()[0];
        if (audioTrack) {
            audioTrack.enabled = !audioTrack.enabled;
            controls.toggleMic.querySelector('.material-icons').textContent =
                audioTrack.enabled ? 'mic' : 'mic_off';
        }
    };

    controls.toggleVideo.onclick = () => {
        const videoTrack = calls.localStream?.getVideoTracks()[0];
        if (videoTrack) {
            videoTrack.enabled = !videoTrack.enabled;
            controls.toggleVideo.querySelector('.material-icons').textContent =
                videoTrack.enabled ? 'videocam' : 'videocam_off';
        }
    };

    controls.shareScreen.onclick = async () => {
        if (!screenShareCalls) {
            await setupScreenShare();
        } else {
            await stopScreenShare();
        }
    };

    controls.leave.onclick = async () => {
        if (currentRoom) {
            await calls.leaveRoom();
            window.location.href = 'index.html';
        }
    };

    controls.toggleMask.onclick = () => {
        if (!isMaskEnabled) {
            showMaskModal();
        } else {
            isMaskEnabled = false;
            const maskBtn = document.getElementById('toggleMaskBtn');
            maskBtn.classList.remove('active');
            combineEffects();
        }
    };

    // Add blur toggle handler
    document.getElementById('toggleBlurBtn').onclick = async () => {
        isBlurEnabled = !isBlurEnabled;
        const blurBtn = document.getElementById('toggleBlurBtn');
        blurBtn.classList.toggle('active', isBlurEnabled);
        await combineEffects();
    };
}

// Thêm CSS styles cho hiệu ứng vẫy tay
const style = document.createElement('style');
style.textContent = `
    .wave-effect {
        position: absolute;
        top: 10px;
        right: 10px;
        font-size: 24px;
        animation: wave 1s infinite;
    }
    
    @keyframes wave {
        0% { transform: rotate(0deg); }
        25% { transform: rotate(-20deg); }
        75% { transform: rotate(20deg); }
        100% { transform: rotate(0deg); }
    }
`;
document.head.appendChild(style);

function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `notification ${type}`;
    notification.textContent = message;
    notificationsContainer.appendChild(notification);
    setTimeout(() => notification.remove(), 3000);
}

// Initialize when page loads
async function initialize() {
    if (!username || !roomId) {
        window.location.href = 'index.html';
        return;
    }

    // Initialize background blur with existing canvas element
    const blurCanvas = document.getElementById('blurCanvas');
    if (!blurCanvas) {
        console.error('Blur canvas element not found');
        return;
    }
    
    backgroundBlur = new BackgroundBlur(
        document.createElement('video'),
        blurCanvas
    );
    await backgroundBlur.initialize();

    // Load available masks
    await loadAvailableMasks();

    // Initialize face mask filter
    const maskCanvas = document.getElementById('maskCanvas');
    const maskImage = document.getElementById('maskImage');
    if (maskCanvas && maskImage) {
        faceMaskFilter = new FaceMaskFilter(
            document.createElement('video'),
            maskCanvas,
            maskImage
        );
        await faceMaskFilter.initialize();
    }

    if (await ensureInitialized()) {
        await setupLocalVideo();
        await joinRoom();
        setupCallbacks();
    }
}

document.addEventListener('DOMContentLoaded', initialize);

window.addEventListener('beforeunload', () => {
    if (currentRoom) {
        calls.leaveRoom();
        if (screenShareCalls) {
            screenShareCalls.leaveRoom();
        }
    }
});

// Add mask-related functions
async function loadAvailableMasks() {
    masksList = [
        'basic/mask1.png',
        'basic/mask2.png',
        'basic/mask3.png',
        'medicel/mask1.png',
        'medicel/mask2.png',
        'medicel/mask3.png'
    ];
    console.log('Available masks:', masksList);
}

// Add modal close button handler
document.querySelector('#maskModal .close-btn').onclick = () => {
    document.getElementById('maskModal').classList.remove('show');
};

function showMaskModal() {
    const modal = document.getElementById('maskModal');
    const maskGrid = document.getElementById('maskGrid');
    maskGrid.innerHTML = '';

    // Group masks by category
    const masksByCategory = masksList.reduce((acc, maskFile) => {
        const category = maskFile.split('/')[0];
        if (!acc[category]) {
            acc[category] = [];
        }
        acc[category].push(maskFile);
        return acc;
    }, {});

    // Create mask options for each category
    Object.entries(masksByCategory).forEach(([category, masks]) => {
        masks.forEach(maskFile => {
            const maskOption = document.createElement('div');
            maskOption.className = `mask-option ${maskFile === currentMask ? 'selected' : ''}`;
            
            const maskName = maskFile.split('/')[1].replace('.png', '');

            maskOption.innerHTML = `
                <img src="assets/mask/${maskFile}" alt="${maskName}">
                <div class="mask-name">${maskName}</div>
                <div class="mask-category">${category}</div>
            `;

            maskOption.onclick = () => {
                document.querySelectorAll('.mask-option').forEach(opt =>
                    opt.classList.remove('selected')
                );
                maskOption.classList.add('selected');
                currentMask = maskFile;
                isMaskEnabled = true;
                updateMaskState();
            };

            maskGrid.appendChild(maskOption);
        });
    });

    // Show modal with animation
    modal.classList.add('show');
    setTimeout(() => modal.querySelector('.modal-content').classList.add('show'), 50);
}

// Update modal close handler
document.querySelector('#maskModal .close-btn').onclick = () => {
    const modal = document.getElementById('maskModal');
    modal.querySelector('.modal-content').classList.remove('show');
    setTimeout(() => modal.classList.remove('show'), 300);
};

// Add click outside to close
document.getElementById('maskModal').onclick = (e) => {
    if (e.target.id === 'maskModal') {
        e.target.querySelector('.modal-content').classList.remove('show');
        setTimeout(() => e.target.classList.remove('show'), 300);
    }
};

async function updateMaskState() {
    const maskBtn = document.getElementById('toggleMaskBtn');
    maskBtn.classList.toggle('active', isMaskEnabled);

    const maskImage = document.getElementById('maskImage');
    maskImage.src = `assets/mask/${currentMask}`;

    document.getElementById('maskModal').classList.remove('show');

    await combineEffects();
}

async function updateBlurState() {
    const blurBtn = document.getElementById('toggleBlurBtn');
    blurBtn.classList.toggle('active', isBlurEnabled);

    await combineEffects();
}

async function combineEffects() {
    try {
        let finalStream = calls.localStream;

        if (isMaskEnabled && isBlurEnabled) {
            // Apply mask first
            const maskedStream = await faceMaskFilter.processFrame(calls.localStream);
            await new Promise(resolve => setTimeout(resolve, 100)); // Wait for mask to initialize

            // Then apply blur
            const blurredAndMaskedStream = await backgroundBlur.updateInputStream(maskedStream);

            // Create final stream with both effects
            finalStream = new MediaStream();
            blurredAndMaskedStream.getVideoTracks().forEach(track => {
                finalStream.addTrack(track);
            });

            // Add audio track
            const audioTrack = calls.localStream.getAudioTracks()[0];
            if (audioTrack) {
                finalStream.addTrack(audioTrack);
            }

        } else if (isMaskEnabled) {
            finalStream = await faceMaskFilter.processFrame(calls.localStream);
            // Ensure we have audio
            const audioTrack = calls.localStream.getAudioTracks()[0];
            if (audioTrack && !finalStream.getAudioTracks().length) {
                finalStream.addTrack(audioTrack);
            }

        } else if (isBlurEnabled) {
            const blurredStream = await backgroundBlur.updateInputStream(calls.localStream);
            finalStream = new MediaStream();
            blurredStream.getVideoTracks().forEach(track => {
                finalStream.addTrack(track);
            });
            
            // Add audio track
            const audioTrack = calls.localStream.getAudioTracks()[0];
            if (audioTrack) {
                finalStream.addTrack(audioTrack);
            }
        }

        // Update local video display với muted audio
        const localVideo = document.querySelector('.local-video video');
        if (localVideo) {
            localVideo.srcObject = finalStream;
            localVideo.muted = true; // Đảm bảo local preview luôn mute
            localVideo.volume = 0;
        }

        // Update WebRTC stream
        if (calls.peerConnection) {
            const videoSender = calls.peerConnection.getSenders()
                .find(sender => sender.track?.kind === 'video');
            if (videoSender) {
                const videoTrack = finalStream.getVideoTracks()[0];
                if (videoTrack) {
                    await videoSender.replaceTrack(videoTrack);
                }
            }
        }

    } catch (error) {
        console.error('Error in combineEffects:', error);
    }
}
