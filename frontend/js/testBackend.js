// filepath: d:/DappMeetingV3/public/js/testBackend.js

async function testTokenGeneration() {
    const response = await fetch('http://localhost:5000/auth/token', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: 'testuser' })
    });
    const data = await response.json();
    console.log('Token Generation:', data);
    return data.token;
}

async function testCreateRoom(token) {
    const response = await fetch('http://localhost:5000/api/rooms', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ name: 'Test Room', metadata: { description: 'A test room' } })
    });
    const data = await response.json();
    console.log('Create Room:', data);
    return data.roomId;
}

async function testJoinRoom(token, roomId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({})
    });
    const data = await response.json();
    console.log('Join Room:', data);
    return data.sessionId;
}

async function testPublishTracks(token, roomId, sessionId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/sessions/${sessionId}/publish`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({
            offer: { sdp: 'dummy sdp', type: 'offer' },
            tracks: [{ trackName: 'video', mid: '0', location: 'local' }]
        })
    });
    const data = await response.json();
    console.log('Publish Tracks:', data);
}

async function testUnpublishTrack(token, roomId, sessionId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/sessions/${sessionId}/unpublish`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({
            trackName: 'video',
            mid: '0',
            force: false,
            sessionDescription: { sdp: 'dummy sdp', type: 'offer' }
        })
    });
    const data = await response.json();
    console.log('Unpublish Track:', data);
}

async function testPullTracks(token, roomId, sessionId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/sessions/${sessionId}/pull`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ remoteSessionId: sessionId, trackName: 'video' })
    });
    const data = await response.json();
    console.log('Pull Tracks:', data);
}

async function testRenegotiateSession(token, roomId, sessionId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/sessions/${sessionId}/renegotiate`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ sessionDescription: { sdp: 'dummy sdp', type: 'offer' } })
    });
    const data = await response.json();
    console.log('Renegotiate Session:', data);
}

async function testManageDataChannels(token, roomId, sessionId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/sessions/${sessionId}/datachannels/new`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ dataChannels: [{ location: 'local', dataChannelName: 'chat' }] })
    });
    const data = await response.json();
    console.log('Manage Data Channels:', data);
}

async function testGetSessionState(token, roomId, sessionId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/sessions/${sessionId}/state`, {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const data = await response.json();
    console.log('Get Session State:', data);
}

async function testGetICEServers(token) {
    const response = await fetch('http://localhost:5000/ice-servers', {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const data = await response.json();
    console.log('Get ICE Servers:', data);
}

async function testGetUserInfo(token) {
    const response = await fetch('http://localhost:5000/api/users/me', {
        method: 'GET',
        headers: { 'Authorization': `Bearer ${token}` }
    });
    const data = await response.json();
    console.log('Get User Info:', data);
}

async function testLeaveRoom(token, roomId, sessionId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/leave`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ sessionId })
    });
    const data = await response.json();
    console.log('Leave Room:', data);
}

async function testUpdateTrackStatus(token, roomId, sessionId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/sessions/${sessionId}/track-status`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ trackId: 'video', kind: 'video', enabled: false, force: false })
    });
    const data = await response.json();
    console.log('Update Track Status:', data);
}

async function testUpdateRoomMetadata(token, roomId) {
    const response = await fetch(`http://localhost:5000/api/rooms/${roomId}/metadata`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ name: 'Updated Room Name', metadata: { description: 'Updated description' } })
    });
    const data = await response.json();
    console.log('Update Room Metadata:', data);
}

async function runTests() {
    const token = await testTokenGeneration();
    const roomId = await testCreateRoom(token);
    const sessionId = await testJoinRoom(token, roomId);
    await testPublishTracks(token, roomId, sessionId);
    await testUnpublishTrack(token, roomId, sessionId);
    await testPullTracks(token, roomId, sessionId);
    await testRenegotiateSession(token, roomId, sessionId);
    await testManageDataChannels(token, roomId, sessionId);
    await testGetSessionState(token, roomId, sessionId);
    await testGetICEServers(token);
    await testGetUserInfo(token);
    await testLeaveRoom(token, roomId, sessionId);
    await testUpdateTrackStatus(token, roomId, sessionId);
    await testUpdateRoomMetadata(token, roomId);
}

runTests();
