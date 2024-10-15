let token = null;

function login() {
  const username = document.getElementById('username').value;
  const password = document.getElementById('password').value;

  fetch('http://localhost:8080/api/v1/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, password }),
  })
    .then(response => response.json())
    .then(data => {
      if (data.access_token) {
        token = data.access_token;
        localStorage.setItem('token', token);
        document.getElementById('chatSection').style.display = 'block';
      } else {
        alert('Login failed');
      }
    });
}

function fetchMessages(receiverId) {
  fetch(`http://localhost:8080/api/v1/messages?sender_id=6&receiver_id=${receiverId}`, {
    headers: { 'Authorization': `Bearer ${localStorage.getItem('token')}` },
  })
    .then(response => response.json())
    .then(data => {
      const messagesDiv = document.getElementById('messages');
      messagesDiv.innerHTML = '';
      data.forEach(msg => {
        const messageDiv = document.createElement('div');
        messageDiv.className = 'message';
        messageDiv.textContent = `${msg.sender_id === 6 ? 'You' : 'Them'}: ${msg.content} (${new Date(msg.sent_at).toLocaleTimeString()})`;
        messagesDiv.appendChild(messageDiv);
      });
    });
}

function sendMessage() {
  const receiverId = document.getElementById('receiver').value;
  const content = document.getElementById('messageInput').value;

  fetch('http://localhost:8080/api/v1/messages', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${localStorage.getItem('token')}`,
    },
    body: JSON.stringify({
      sender_id: 6, // Replace with logged-in user ID
      receiver_id: parseInt(receiverId),
      content: content,
    }),
  }).then(() => {
    document.getElementById('messageInput').value = '';
    fetchMessages(receiverId);
  });
}
