document.getElementById('fillForm').addEventListener('submit', async function(e) {
    e.preventDefault();
    document.getElementById("submitURL").style.backgroundColor = "gray"
    document.getElementById("submitURL").disabled = true;
    document.getElementById("submitURL").style.cursor = "default";
    
    const url = document.getElementById('url').value;
    const errorDiv = document.getElementById('error');
    errorDiv.textContent = '';
    
    try {
        const response = await fetch('/fill', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ url: url })
        });
        
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        
        const data = await response.json();
        if (data.status == "fail") {
          return
        }
        
        const msg = JSON.parse(data.msg);
        console.log(msg)
        // Fill in the response fields
        document.getElementById('title').value = msg.title || '';
        document.getElementById('description').value = msg.description || '';

        console.log(msg.img, msg.img.includes(","))
        if (!msg.img.includes(",")) document.getElementById('img1').src = msg.img
        else {
          // fill images. max of 4 images.
          let parts = msg.img.split(",")
          for (let i = 1; i < parts.length && i <= 4; i++) {
            console.log("id:", "img"+i, "src:", parts[i-1])
            console.log('img'+i, document.getElementById('img'+i))
            document.getElementById('img'+i).src = parts[i-1]
          }
        }
        
    } catch (error) {
        errorDiv.textContent = `Error: ${error.message}`;
        console.error('Error:', error);
    }

    document.getElementById("submitURL").style.backgroundColor = "#4CAF50"
    document.getElementById("submitURL").disabled = false;
    document.getElementById("submitURL").style.cursor = "pointer";

    modalOverlay.classList.remove('active');
});


// Magic button
const magicBtn = document.getElementById('magicBtn');
const modalOverlay = document.getElementById('modalOverlay');
const closeBtn = document.getElementById('closeBtn');

// Open modal
magicBtn.addEventListener('click', () => {
  modalOverlay.classList.add('active');
});

// Close modal
closeBtn.addEventListener('click', () => {
  modalOverlay.classList.remove('active');
});

// Close modal when clicking outside
modalOverlay.addEventListener('click', (e) => {
  if (e.target === modalOverlay) {
    modalOverlay.classList.remove('active');
  }
});

// Close modal on Escape key
document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape' && modalOverlay.classList.contains('active')) {
    modalOverlay.classList.remove('active');
  }
});
