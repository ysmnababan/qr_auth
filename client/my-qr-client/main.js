import Pusher from "pusher-js";

let hasQRReceivedOnce = false

// Pusher initialization
const clientToken = crypto.randomUUID();  // Generates a UUID v4
const pusher = new Pusher(import.meta.env.VITE_PUSHER_KEY, {
    cluster: import.meta.env.VITE_PUSHER_CLUSTER,
});

if (!hasQRReceivedOnce){
    pusher.connection.bind('connected', function() {
        console.log("Pusher connected!");
    });

    pusher.connection.bind('error', function(error) {
        console.log("Pusher error:", error);
    });

    // Subscribe to the channel
    const channel = pusher.subscribe("qr_channel");
    console.log("Subscribed to qr_channel");
    // Listen for the qr_event and display the QR code
    channel.bind(`qr_event:${clientToken}`, (eventData) => {
        console.log("QR Event Received:", eventData);

        // Display the QR code received from the backend
        const container = document.getElementById("qr-container");
        container.innerHTML = "";  // Clear previous content

        // Create the QR image element
        const img = document.createElement("img");
        img.src = eventData.qr_code;  // Use the qr_code from the Pusher event
        img.alt = "QR Code";
        img.width = 256;
        img.height = 256;

        container.appendChild(img);  // Append to the container
    });

    // Listen for the success event (or poll for URL visit) to hide the QR code
    channel.bind(`login_success:${clientToken}`, (eventData) => {
        console.log("Login successful, hiding QR code...", eventData.status);
        // Hide the QR code container
        const container = document.getElementById("qr-container");
        container.innerHTML = "";  // Clear the QR code

        // Optionally show a success message
        
        const messageContainer = document.getElementById("login-message");
        messageContainer.innerHTML = "<p>Login successful!</p>";  // You can add more styling here
        if (eventData.status === "true"){
            console.log(eventData.token)

             // 👉 Unsubscribe from channel
            pusher.unsubscribe("qr_channel");

            // 👉 Disconnect Pusher connection
            pusher.disconnect();

            pusher.connection.bind('connected', function() {
                console.log("Pusher connected!");
            });
            setTimeout(() => {
                window.location.href = "/dashboard.html";  // adjust path based on your setup
            }, 100);
        }
    });
}

// Event listener for "Start Login" button
document.getElementById("start-login").addEventListener("click", async () => {
    try {
        // Hide the button when clicked
        const button = document.getElementById("start-login");
        button.style.display = "none";  // Hides the button
        
        // Trigger the backend login process
        const baseURL = "https://d13e-125-161-205-3.ngrok-free.app"
        const url = `${baseURL}/auth/qr-login?uuid=${encodeURIComponent(clientToken)}`;
        console.log(url)
        const response = await fetch(url, {
            method: "GET",
            headers: {
                'ngrok-skip-browser-warning': 'true', // This skips the ngrok warning page
            },
        });
        // const text = await response.text();  // ← Use this for now
        // console.log("Raw response:", text);

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
    
        // Parse the JSON response (assuming backend sends QR code as part of response)
        const data = await response.json();  // Parse the JSON response
        console.log("QR Data received from backend:", data);
        hasQRReceivedOnce = true
    }catch (error) {
        alert("Error starting login process: " + error.message);
        console.error("Fetch error:", error);
    }
});