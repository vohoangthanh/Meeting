// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract DAppMeeting {
    // Structures
    struct Track {
        string trackName;   // Name of the track
        string mid;         // Media identifier in WebRTC
        string location;    // Track location
        bool isPublished;   // Publication status
        string sessionId;   // Session ID associated with this track
        string roomId;      // Room ID

    }

    struct Participant {
        address walletAddress;  // Participant's wallet address
        string name;            // Participant's name
        string sessionID;       // Session ID from Cloudflare Calls
    }

    struct Room {
        string roomId;          // Room ID
        uint256 creationTime;   // Creation timestamp
        Participant[] participants; // List of participants
    }

    // State variables
    mapping(string => Room) public rooms; // Store rooms by roomId
    mapping(string => mapping(address => uint)) public participantIndices; // Track participant indices in arrays
    mapping(string => mapping(address => bool)) public participantsInRoom; // Check if participant is in room
    
    // New track mapping: roomId -> participantAddress -> trackList
    mapping(string => mapping(address => Track[])) public participantTracks;
    // Track count mapping for easy querying
    mapping(string => mapping(address => uint)) public participantTrackCount;
    
    address public owner;
    address[] public authorizedBackends;
    
    // Events
    event ParticipantJoined(string roomId, address participant, Track[] initialTracks, bytes sessionDescription);
    event ParticipantLeft(string roomId, address participant);
    event TrackAdded(string roomId, address participant, string trackName);
    event EventForwardedToBackend(string roomId, address sender, bytes eventData);
    event EventForwardedToFrontend(string roomId, address participant, bytes eventData);
    
    // Constructor
    constructor() {
        owner = msg.sender;
        authorizedBackends.push(msg.sender); // Owner is authorized by default
    }
    
    // Modifiers
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }
    
    modifier onlyAuthorized() {
        bool isAuthorized = false;
        for (uint i = 0; i < authorizedBackends.length; i++) {
            if (msg.sender == authorizedBackends[i]) {
                isAuthorized = true;
                break;
            }
        }
        require(isAuthorized, "Not authorized");
        _;
    }
    
    modifier roomExists(string memory _roomId) {
        require(bytes(rooms[_roomId].roomId).length > 0, "Room does not exist");
        _;
    }
    
    modifier participantExists(string memory _roomId) {
        require(participantsInRoom[_roomId][msg.sender], "You are not in this room");
        _;
    }
    
    // Management functions
    function addAuthorizedBackend(address _backend) public {
        authorizedBackends.push(_backend);
    }
    
    function removeAuthorizedBackend(address _backend) public {
        for (uint i = 0; i < authorizedBackends.length; i++) {
            if (authorizedBackends[i] == _backend) {
                // Swap with the last element and remove
                authorizedBackends[i] = authorizedBackends[authorizedBackends.length - 1];
                authorizedBackends.pop();
                break;
            }
        }
    }
    
    // Core functions
    function createRoom(string memory _roomId) public {
        require(bytes(rooms[_roomId].roomId).length == 0, "Room already exists");
        
        // Create a new room
        rooms[_roomId].roomId = _roomId;
        rooms[_roomId].creationTime = block.timestamp;
    }
    
    function joinRoom(string memory _roomId, string memory _name, Track[] memory _initialTracks, bytes memory sessionDescription) public roomExists(_roomId) {
        require(!participantsInRoom[_roomId][msg.sender], "Already in room");
        
        // Add participant to room
        Participant memory newParticipant = Participant({
            walletAddress: msg.sender,
            name: _name,
            sessionID: ""
        });
        
        rooms[_roomId].participants.push(newParticipant);
        uint participantIndex = rooms[_roomId].participants.length - 1;
        participantIndices[_roomId][msg.sender] = participantIndex;
        participantsInRoom[_roomId][msg.sender] = true;
        
        // Add initial tracks
        for (uint i = 0; i < _initialTracks.length; i++) {
            Track memory track = _initialTracks[i];
            participantTracks[_roomId][msg.sender].push(track);
            participantTrackCount[_roomId][msg.sender]++;
            emit TrackAdded(_roomId, msg.sender, track.trackName);
        }
        
        // Emit the event for backend
        emit ParticipantJoined(_roomId, msg.sender, _initialTracks, sessionDescription);
    }
    
    function leaveRoom(string memory _roomId) public roomExists(_roomId) participantExists(_roomId) {
        uint participantIndex = participantIndices[_roomId][msg.sender];
        
        // Delete participant from the room
        // Move the last participant to the position of the leaving participant
        if (participantIndex < rooms[_roomId].participants.length - 1) {
            Participant memory lastParticipant = rooms[_roomId].participants[rooms[_roomId].participants.length - 1];
            rooms[_roomId].participants[participantIndex] = lastParticipant;
            participantIndices[_roomId][lastParticipant.walletAddress] = participantIndex;
        }
        
        // Remove the last element
        rooms[_roomId].participants.pop();
        
        // Update mappings
        participantsInRoom[_roomId][msg.sender] = false;
        delete participantIndices[_roomId][msg.sender];
        delete participantTracks[_roomId][msg.sender];
        delete participantTrackCount[_roomId][msg.sender];
        
        // Emit event
        emit ParticipantLeft(_roomId, msg.sender);
    }
    
    function setParticipantSessionID(string memory _roomId, address _participantAddress, string memory _sessionID) 
        public roomExists(_roomId) {
        require(participantsInRoom[_roomId][_participantAddress], "Participant not in room");
        
        uint participantIndex = participantIndices[_roomId][_participantAddress];
        rooms[_roomId].participants[participantIndex].sessionID = _sessionID;
    }
    
    function addTrack(string memory _roomId, Track memory _newTrack) public roomExists(_roomId) participantExists(_roomId) {
        participantTracks[_roomId][msg.sender].push(_newTrack);
        participantTrackCount[_roomId][msg.sender]++;
        emit TrackAdded(_roomId, msg.sender, _newTrack.trackName);
    }
    
    function getParticipantInfo(string memory _roomId) public view roomExists(_roomId) participantExists(_roomId) returns (Participant memory) {
        uint participantIndex = participantIndices[_roomId][msg.sender];
        return rooms[_roomId].participants[participantIndex];
    }
    
    function forwardEventToBackend(string memory _roomId, bytes memory _eventData) public roomExists(_roomId) participantExists(_roomId) {
        emit EventForwardedToBackend(_roomId, msg.sender, _eventData);
    }
    
    function forwardEventToFrontend(string memory _roomId, address _participant, bytes memory _eventData) public roomExists(_roomId) {
        require(participantsInRoom[_roomId][_participant], "Target participant not in room");
        emit EventForwardedToFrontend(_roomId, _participant, _eventData);
    }
    
    // New function to get a participant's tracks
    function getParticipantTracks(string memory _roomId, address _participant) public view roomExists(_roomId) returns (Track[] memory) {
        require(participantsInRoom[_roomId][_participant], "Participant not in room");
        return participantTracks[_roomId][_participant];
    }
    
    // Helper functions
    function getRoomParticipantsCount(string memory _roomId) public view roomExists(_roomId) returns (uint) {
        return rooms[_roomId].participants.length;
    }
    
    function getParticipantTracksCount(string memory _roomId, address _participant) public view roomExists(_roomId) returns (uint) {
        require(participantsInRoom[_roomId][_participant], "Participant not in room");
        return participantTrackCount[_roomId][_participant];
    }
}