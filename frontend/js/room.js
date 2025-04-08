import CloudflareCalls from './CloudflareCalls.js';
import SmartContractConnector from './smartContractIntegration.js';
import { auth } from './auth.js';
import eventListener from './contractEventListener.js';

import smartContractIntegration from './smartContractIntegration.js';

// Lấy current URL và cấu hình endpoints
// const currentUrl = 'https://manhteky123-dappmeetingv3.hf.space';
const currentUrl = window.location.origin;
const isLocalhost = currentUrl.includes('localhost') || currentUrl.includes('127.0.0.1');

// Cấu hình endpoints dựa trên môi trường
const config = {
    backendUrl: isLocalhost ? 'http://localhost:50000' : currentUrl,
    websocketUrl: isLocalhost
        ? 'ws://localhost:50000/ws'
        : 'wss://manhteky123-dappmeetingv3.hf.space/ws'
};
const baseAPI = isLocalhost ? 'http://localhost:50000' : 'https://manhteky123-dappmeetingv3.hf.space';
const calls = new CloudflareCalls(config);

// Tương tự cho screenShareCalls
const screenShareConfig = { ...config };

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
let participants = [];

let trackPullCompleteListener = null;
let leaveRoomListener = null;

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
const walletAddress = localStorage.getItem('privateKey');


// Check wallet connection before proceeding
async function checkWalletAuth() {
    // Check if wallet address exists in localStorage
    if (!walletAddress) {
        showNotification('Wallet authentication required', 'error');
        setTimeout(() => {
            window.location.href = 'index.html';
        }, 2000);
        return false;
    }

    return true;
}

// Get token and initialize calls
async function ensureInitialized() {
    // if (!calls.token) {
    //     try {
    //         const response = await fetch(`${baseAPI}/auth/token`, {
    //             method: 'POST',
    //             headers: { 'Content-Type': 'application/json' },
    //             body: JSON.stringify({
    //                 username,
    //                 walletAddress: auth.userAddress
    //             })
    //         });

    //         const { token } = await response.json();
    //         calls.setToken(token);
    //         showNotification('Successfully initialized');
    //         return true;
    //     } catch (err) {
    //         console.error('Error getting token:', err);
    //         showNotification('Failed to initialize', 'error');
    //         return false;
    //     }
    // }
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
    let displayName = participant.name || `User-${participant.userId.slice(0, 6)}`;
    if (participant.walletAddress) {
        const shortAddress = participant.walletAddress.substring(0, 6) + '...' +
            participant.walletAddress.substring(participant.walletAddress.length - 4);
        displayName += ` (${shortAddress})`;
    }
    return displayName;
}

async function pullParticipantTracks(sessionId, trackNames) {
    try {
        for (const trackName of trackNames) {
            console.log(`Pulling track ${trackName} from session ${sessionId}`);
            // Process one track at a time, waiting for each to complete
            await calls._pullTracks(sessionId, trackName);
            console.log(`Successfully pulled track ${trackName} from session ${sessionId}`);
        }
        return true;
    } catch (error) {
        console.error('Error pulling tracks:', error);
        return false;
    }
}

async function joinRoom() {
    // 1. First verify smart contract connection is established
    if (!SmartContractConnector.contract) {
        console.error('Smart contract not initialized, initializing...');
        await SmartContractConnector.initialize();
    }

    // Thiết lập listener để nhận answer từ backend và establish peer connection
    const removeAnswerListener = await SmartContractConnector.listenForTrackPublishedAnswer(
        roomId,
        async (answerData) => {
            try {
                console.log('Received track publish answer from backend:', answerData);
                
                if (answerData && answerData.sessionDescription) {
                    // Tạo và thiết lập remote description từ answer
                    if (calls.peerConnection && calls.peerConnection.signalingState === 'have-local-offer') {
                        console.log('Setting remote description from backend answer');
                        const remoteDesc = new RTCSessionDescription(answerData.sessionDescription);
                        await calls.peerConnection.setRemoteDescription(remoteDesc);
                        console.log('Remote description set successfully, peer connection established');
                    } else if (calls.peerConnection) {
                        console.warn(`Invalid signaling state for setRemoteDescription: ${calls.peerConnection.signalingState}`);
                    } else {
                        console.warn('PeerConnection not available');
                    }
                }
            } catch (error) {
                console.error('Error processing track publish answer:', error);
            }
        }
    );

    // 2. Before joining the room via WebRTC, make sure to join via smart contract first
    // This ensures you're properly registered as a participant
    console.log('Joining room via smart contract first...');
    await SmartContractConnector.joinRoom(roomId, username, []);
    console.log('Successfully joined room via smart contract');
    
    // 1. Join room và lấy session
    await calls.joinRoom(roomId, { name: username });
    currentRoom = roomId;
    showNotification(`Joined room: ${roomId}`);


    // Đảm bảo event listener đã được khởi tạo
    if (!eventListener) {
        console.error('Event listener not initialized');
        throw new Error('Event listener not initialized');
    }


    // Start listening for events from backend before proceeding
    console.log('Starting to listen for backend events');
    

    // 2. Setup handlers trước khi pull tracks
    setupCallbacks();

    // Verify smart contract connection
    if (!SmartContractConnector.contract) {
        console.error('Smart contract not initialized, reinitializing...');
        await SmartContractConnector.initialize();
    }

    // 3. Lấy danh sách người tham gia
    participants = await calls.getParticipantsFromContract(roomId);
    console.log('Current participants:', participants);

    // 4. Cập nhật sessionId của bản thân khi lấy từ smart constract
    try {
        // Create wallet from private key to get the address
        const wallet = new ethers.Wallet(localStorage.getItem('privateKey'));
        const myWalletAddress = wallet.address;
        
        console.log('Looking for participant with my wallet address:', myWalletAddress);
        
        // Find participant with matching wallet address
        const participant = participants.find(p => 
            p.walletAddress && 
            p.walletAddress.toLowerCase() === myWalletAddress.toLowerCase()
        );
        
        if (participant) {
            calls.sessionId = participant.sessionId;
            console.log('Found myself in participants list. My session ID:', calls.sessionId);
        } else {
            console.error('Participant not found in the list. Available participants:', 
                participants.map(p => ({ name: p.name, address: p.walletAddress, sessionId: p.sessionId }))
            );
            return;
        }
    } catch (error) {
        console.error('Error finding participant:', error);
        return;
    }

    // Đợi peerConnection kết nối thành công
    await new Promise(resolve => {
        console.log('Peer Connection state:', calls.peerConnection.connectionState);
        console.log('Waiting for connection to be established...');
        if (calls.peerConnection.connectionState === 'connected') {
            resolve();
        } else {
            calls.peerConnection.onconnectionstatechange = () => {
                if (calls.peerConnection.connectionState === 'connected') {
                    resolve();
                }
            };
        }
    });

    console.log('Peer Connection established:', calls.peerConnection.connectionState);

    if(calls.sessionId === null || calls.sessionId === undefined || calls.sessionId === '') {
        console.error('Session ID is null or undefined. Cannot proceed.');
        return;

    }
    var userAddress = localStorage.getItem('wallet_address');
    // 5. Set up remote streams for existing participants
    for (const participant of participants) {
        // Skip if it's our own session
        if (participant.walletAddress.toLowerCase() === userAddress.toLowerCase()) continue;

        console.log('Processing participant:', participant);

        // Pull each track from the participant
        for (const trackName of participant.publishedTracks) {
            console.log(`Pulling track ${trackName} from session ${participant.sessionId}`);
            await calls._pullTracks(participant.sessionId, trackName);
        }
    }

    // Set session ID cho contract event listener
    eventListener.setCurrentSessionId(calls.sessionId);

    // 5. Start monitoring stats sau khi mọi thứ đã setup
    calls.startStatsMonitoring(1000);
}


async function checkAndRestoreEventListeners() {
    // const reconnected = await SmartContractConnector.reestablishEventListeners();
    const reconnected = true;

    if (reconnected) {
        console.log("WebSocket reconnected, restoring event listeners...");

        // Re-establish track pull complete listener
        if (trackPullCompleteListener) {
            trackPullCompleteListener = await SmartContractConnector.listenForTrackPullComplete(
                roomId,
                calls.sessionId,
                async (data) => {
                    console.log('Track pull complete:', data);
                    if (data.sessionId !== calls.sessionId) {
                        console.log(`New participant joined with ${data.trackNames.length} tracks, pulling tracks...`);
                        await pullParticipantTracks(data.sessionId, data.trackNames);
                    }
                }
            );
        }

        // Re-establish leave room listener
        if (leaveRoomListener) {
            leaveRoomListener = await SmartContractConnector.listenForParticipantLeaveRoom(
                roomId,
                (data) => {
                    console.log('Participant left:', data);
                    try {
                        const sessionId = data.sessionId;
                        const container = document.getElementById(`participant-${sessionId}`);
                        if (container) {
                            container.remove();
                        }
                        const container2 = document.getElementById(`participant ${sessionId}`);
                        if (container2) {
                            container2.remove();
                        }
                        showNotification(`A participant left the room`);
                    } catch (error) {
                        console.error('Error handling participant leave event:', error);
                    }
                }
            );
        }
    }
}

async function setupScreenShare() {
    try {
        // Get screen share stream
        const screenStream = await navigator.mediaDevices.getDisplayMedia({
            video: true,
            audio: false
        });

        // Store screen stream reference
        calls.screenStream = screenStream;

        // Create new stream for screen sharing
        const screenTrack = screenStream.getVideoTracks()[0];

        // Create new peer connection for screen share if needed
        if (!calls.peerConnection) {
            await calls._createPeerConnection();
        }


        // Đợi trạng thái kết nối ổn định trước khi thêm track
        if (calls.peerConnection.signalingState !== 'stable') {
            console.log('Waiting for signaling state to stabilize...');
            await new Promise(resolve => {
                const checkState = () => {
                    if (calls.peerConnection.signalingState === 'stable') {
                        resolve();
                    } else {
                        setTimeout(checkState, 100);
                    }
                };
                checkState();
            });
        }
        // First add the track to the peer connection and store the sender
        calls.screenTrackSender = calls.peerConnection.addTrack(screenTrack, calls.localStream);


        calls.localStream.addTrack(screenTrack);

        // Create screen share container
        const screenContainer = document.createElement('div');
        screenContainer.id = 'screen-share-container';
        screenContainer.className = 'video-container screen-share';

        // Create video element for screen share
        const screenVideo = document.createElement('video');
        screenVideo.autoplay = true;
        screenVideo.playsInline = true;
        screenVideo.srcObject = calls.screenStream;

        // Add label
        const nameLabel = document.createElement('div');
        nameLabel.className = 'participant-name';
        nameLabel.textContent = 'Screen Share';

        screenContainer.appendChild(screenVideo);
        screenContainer.appendChild(nameLabel);
        videoGrid.appendChild(screenContainer);

        // Publish track
        try {
            await calls._publishTracks();
            console.log("Screen share track published successfully");
            const wallet = new ethers.Wallet(localStorage.getItem('privateKey'));
            // Get user address for event filtering
            const userAddress = wallet.address;
            // Get your published tracks
            console.log("Notify share track: ", roomId, userAddress, calls.sessionId, [screenTrack.id]);
            await smartContractIntegration.notifyTrackPullComplete(
                roomId,
                userAddress,
                calls.sessionId,
                [screenTrack.id]
            );
            // Check and restore WebSocket event listeners
            await checkAndRestoreEventListeners();
        } catch (error) {
            console.error("Error publishing screen share track:", error);
            screenTrack.stop();
            // Clean up the track if publishing fails
            if (calls.screenTrackSender) {
                calls.peerConnection.removeTrack(calls.screenTrackSender);
                calls.screenTrackSender = null;
            }
            throw error;
        }

        // Listen for end of screen sharing
        screenTrack.onended = async () => {
            await stopScreenShare();
        };

        // Update UI
        controls.shareScreen.querySelector('.material-icons').textContent = 'stop_screen_share';
        showNotification('Screen sharing started');

    } catch (error) {
        console.error("Error setting up screen share:", error);
        showNotification('Failed to share screen', 'error');
    }
}

async function stopScreenShare() {
    try {
        if (calls.screenStream) {
            // Stop all tracks in the screen stream
            calls.screenStream.getTracks().forEach(track => track.stop());

            // Remove the track sender from the peer connection
            if (calls.screenTrackSender) {
                calls.peerConnection.removeTrack(calls.screenTrackSender);
                calls.screenTrackSender = null;
            }

            // Clean up references
            calls.screenStream = null;

            // Update UI
            controls.shareScreen.querySelector('.material-icons').textContent = 'screen_share';
            showNotification('Screen sharing ended');
        }
    } catch (error) {
        console.error("Error stopping screen share:", error);
        showNotification('Error stopping screen share', 'error');
    }
}

async function setupCallbacks() {
    calls.onRemoteTrack((track) => {
        console.log('Remote track received:', track);

        // Create unique ID for each video track
        const trackId = `track-${track.sessionId}-${track.id}`;
        let trackContainer = document.getElementById(trackId);

        // If this track already has a container, just update it
        if (trackContainer) {
            console.log('Container for this track already exists, updating');
            return;
        }

        // For audio tracks, add to the participant's main container
        if (track.kind === 'audio') {
            const participantId = `participant-${track.sessionId}`;
            const participantContainer = document.getElementById(participantId);

            if (participantContainer) {
                const video = participantContainer.querySelector('video');
                if (video && video.srcObject) {
                    video.srcObject.addTrack(track);
                }
            }
            return;
        }

        // For video tracks, always create a new container with unique ID
        trackContainer = document.createElement('div');
        trackContainer.id = trackId;
        trackContainer.className = 'video-container';
        trackContainer.dataset.sessionId = track.sessionId; // Store session ID for cleanup

        const video = document.createElement('video');
        video.autoplay = true;
        video.playsInline = true;

        const name = document.createElement('div');
        name.className = 'participant-name';
        name.textContent = getParticipantName(track.sessionId);

        // Create new MediaStream for this video
        video.srcObject = new MediaStream();
        video.srcObject.addTrack(track);

        trackContainer.appendChild(video);
        trackContainer.appendChild(name);
        videoGrid.appendChild(trackContainer);

        // Helper function to get participant name
        function getParticipantName(sessionId) {
            const participant = participants.find(p => p.sessionId === sessionId);
            return participant ? participant.name : `Participant ${sessionId}`;
        }
    });

    calls.onRemoteTrackUnpublished((sessionId, trackName) => {
        console.log('Remote track unpublished:', { sessionId, trackName });

        // Remove track-specific container if it exists
        const trackContainerId = `track-${sessionId}-${trackName}`;
        const trackContainer = document.getElementById(trackContainerId);
        if (trackContainer) {
            console.log(`Removing track container: ${trackContainerId}`);
            trackContainer.remove();
        }

        // Also check if track exists in participant's container
        const participantContainer = document.getElementById(`participant-${sessionId}`);
        if (participantContainer) {
            const video = participantContainer.querySelector('video');
            if (video && video.srcObject) {
                const stream = video.srcObject;
                const tracks = stream.getTracks();
                tracks.forEach(track => {
                    if (track.id === trackName) {
                        stream.removeTrack(track);
                        track.stop();
                    }
                });
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
        console.log('Participant left:', participant.sessionId);

        // Remove all containers belonging to this participant
        const containers = document.querySelectorAll(`[data-session-id="${participant.sessionId}"]`);
        containers.forEach(container => {
            console.log(`Removing container for left participant: ${container.id}`);
            container.remove();
        });

        // Also remove traditional containers
        const participantContainer = document.getElementById(`participant-${participant.sessionId}`);
        if (participantContainer) {
            participantContainer.remove();
        }

        // Remove any containers with old format IDs
        const oldContainer = document.getElementById(`participant ${participant.sessionId}`);
        if (oldContainer) {
            oldContainer.remove();
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
            await SmartContractConnector.sendNotificationToRoom(roomId, '👋 ' + username + ' vẫy tay chào!');
            console.log('Wave sent successfully');
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
            try {
                // Hiển thị loading indicator hoặc thông báo
                showNotification('Leaving room, please wait...', 'info');
                controls.leave.disabled = true;  // Disable nút để tránh click nhiều lần

                // Đợi cho đến khi leaveRoom hoàn tất
                await calls.leaveRoom();

                // Nếu thành công, chuyển trang
                showNotification('Successfully left room', 'success');
                window.location.href = 'index.html';
            } catch (error) {
                console.error('Error leaving room:', error);
                showNotification('Failed to leave room properly', 'error');
                controls.leave.disabled = false;  // Re-enable nút nếu có lỗi
            }
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


    try {
        // Setup event listeners concurrently using Promise.all
        const [trackPullCleanup, leaveRoomCleanup] = await Promise.all([
            // Track pull complete listener 
            SmartContractConnector.listenForTrackPullComplete(roomId, calls.sessionId, async (data) => {
                console.log('Track pull complete:', data);
                if (data.sessionId !== calls.sessionId) {
                    console.log(`New participant joined with ${data.trackNames.length} tracks, pulling tracks...`);
                    await pullParticipantTracks(data.sessionId, data.trackNames);
                }
            }),

            // Leave room listener
            SmartContractConnector.listenForParticipantLeaveRoom(roomId, (data) => {
                console.log('Participant left:', data);
                try {
                    const sessionId = data.sessionId;
                    const container = document.getElementById(`participant-${sessionId}`);
                    if (container) {
                        container.remove();
                    }
                    const container2 = document.getElementById(`participant ${sessionId}`);
                    if (container2) {
                        container2.remove();
                    }
                    showNotification(`A participant left the room`);
                } catch (error) {
                    console.error('Error handling participant leave event:', error);
                }
            })
        ]);

        // Save listeners to global variables
        trackPullCompleteListener = trackPullCleanup;
        leaveRoomListener = leaveRoomCleanup;

    } catch (error) {
        console.error('Error setting up event listeners:', error);
        throw error;
    }
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

    // Check wallet authentication first
    if (!await checkWalletAuth()) {
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
    }
}

document.addEventListener('DOMContentLoaded', initialize);

window.addEventListener('beforeunload', () => {
    if (currentRoom) {
        // If we have an active screen share, stop it first
        if (calls.screenStream) {
            calls.screenStream.getTracks().forEach(track => track.stop());
        }
        calls.leaveRoom();
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

document.addEventListener('DOMContentLoaded', initialize);

// Listen for wallet disconnection events
auth.addAccountsChangedListener((address) => {
    if (!address) {
        showNotification('Wallet disconnected. You will be redirected to the login page.', 'error');
        setTimeout(() => {
            window.location.href = 'index.html';
        }, 2000);
    }
});
