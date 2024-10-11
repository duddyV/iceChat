const ws = new WebSocket(`ws://${window.location.host}/ws`);
let userField = document.getElementById("username");
let msgField = document.getElementById("message");

document.addEventListener("DOMContentLoaded", function () {
    // ws = new WebSocket("ws://127.0.0.1:8080/ws");
    // const ws = new WebSocket(`ws://${window.location.host}/ws`);

    // WebSocket connection label
    const offline = '<span class="badge bg-danger">Not connected</span>';
    const online = '<span class="badge bg-success">Connected</span>';
    let statusDiv = document.getElementById("status");

    ws.onopen = () => {
        console.log("Succefully connected")
        statusDiv.innerHTML = online;
    }

    ws.onclose = () => {
        console.log("Connection closed");
        statusDiv.innerHTML = offline;
    }

    ws.onerror = error => {
        console.log("ws.onerror issue");
    }

    ws.onmessage = msg => {
        const data = JSON.parse(msg.data);

        if (data.online_users) {
            console.log("Parsed message data", msg.data);
            updateOnlineUsers(data.online_users);
        } else {
            console.log("Message from server ", msg.data);
            document.getElementById("chatbox").innerHTML += `<p><strong>${data.username}:</strong> ${data.message}</p>`;
        }
    }

    // Send username when it changes
    userField.addEventListener("change", () => {
        const newUsername = userField.value || "Anonymous";
        ws.send(JSON.stringify({
            type: "update-username",
            username: newUsername
        }));
    });

    // Validate and send message via btn
    document.getElementById("sendBtn").addEventListener("click", () => {
        if ((userField.value === "") || (msgField.value === "")) {
            errorMsg("Fill out username and messsage!");
            return false;
        } else {
            sendMessage();
        }
    })

    // Handle the Enter button keyup event
    msgField.addEventListener("keyup", function (event) {
        if (event.code === "Enter") {
            if (!ws) {
                console.log("No websocket connection")
                return false
            }

            if ((userField.value === "") || (msgField.value === "")) {
                errorMsg("Fill out username and messsage!");
                return false;
            } else {
                sendMessage();
            }

            event.preventDefault();
            event.stopPropagation();
        }
    })
})

// Update Who's Online list
function updateOnlineUsers(users) {
    const usersList = document.getElementById("users-online");
    usersList.innerHTML = "";
    users.forEach(function (user) {
        const li = document.createElement("li");
        li.textContent = user;
        usersList.appendChild(li);
    });
}

// Handle send message
function sendMessage() {
    let jsonData = {};
    jsonData = {
        type: "chat",
        username: userField.value,
        message: msgField.value,
    }
    ws.send(JSON.stringify(jsonData));
    msgField.value = "";
}

// Notie pretty errors
function errorMsg(msg) {
    notie.alert({
        type: 'error',
        text: msg,
    })
}