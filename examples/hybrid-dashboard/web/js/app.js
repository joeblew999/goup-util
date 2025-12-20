// Hybrid Dashboard - JavaScript
// Communicates with Go backend via HTTP API

// Update stats every second
let updateInterval;

async function updateStats() {
    try {
        const response = await fetch('/api/stats');
        const stats = await response.json();
        
        // Update system info
        document.getElementById('platform').textContent = stats.platform || 'Unknown';
        document.getElementById('goVersion').textContent = stats.goVersion;
        document.getElementById('uptime').textContent = formatUptime(stats.uptime);
        
        // Update CPU gauge
        const cpuPercent = Math.round(stats.cpuUsage);
        document.getElementById('cpuGauge').style.width = cpuPercent + '%';
        document.getElementById('cpuValue').textContent = cpuPercent + '%';
        
        // Update Memory gauge
        const memPercent = Math.round(stats.memoryUsage);
        document.getElementById('memoryGauge').style.width = memPercent + '%';
        document.getElementById('memoryValue').textContent = memPercent + '%';
        
    } catch (error) {
        console.error('Failed to fetch stats:', error);
    }
}

function formatUptime(milliseconds) {
    const seconds = Math.floor(milliseconds / 1000);
    const minutes = Math.floor(seconds / 60);
    const hours = Math.floor(minutes / 60);
    
    if (hours > 0) {
        return hours + 'h ' + (minutes % 60) + 'm';
    } else if (minutes > 0) {
        return minutes + 'm ' + (seconds % 60) + 's';
    } else {
        return seconds + 's';
    }
}

// Call Go function via HTTP
async function callGoFunction() {
    const btn = document.getElementById('callGoBtn');
    const responseDiv = document.getElementById('goResponse');

    btn.disabled = true;
    btn.textContent = 'Calling Go...';

    try {
        const response = await fetch('/api/hello');
        const data = await response.json();

        responseDiv.innerHTML = '<strong>Response from Go:</strong><br>' + JSON.stringify(data, null, 2);

        // Add success animation
        responseDiv.style.background = '#d4edda';
        setTimeout(function() {
            responseDiv.style.background = '#f8f9fa';
        }, 1000);

    } catch (error) {
        responseDiv.innerHTML = '<strong style="color: red;">Error:</strong><br>' + error.message;
    } finally {
        btn.disabled = false;
        btn.textContent = 'Call Go Function';
    }
}

// Check for deep link from Go backend
async function checkDeepLink() {
    try {
        const response = await fetch('/api/deeplink');
        const data = await response.json();

        if (data.hasDeepLink && data.deepLink) {
            const link = data.deepLink;
            document.getElementById('lastDeepLink').textContent = link.originalUrl;
            document.getElementById('deepLinkScheme').textContent = link.scheme;
            document.getElementById('deepLinkPath').textContent = link.path || '/';

            // Highlight the deep link card
            const cards = document.querySelectorAll('.card');
            cards.forEach(card => {
                if (card.querySelector('h2')?.textContent.includes('Deep Linking')) {
                    card.style.borderColor = '#28a745';
                    card.style.boxShadow = '0 0 10px rgba(40, 167, 69, 0.3)';
                    setTimeout(() => {
                        card.style.borderColor = '';
                        card.style.boxShadow = '';
                    }, 3000);
                }
            });

            console.log('Deep link received:', link);
        } else {
            document.getElementById('lastDeepLink').textContent = data.message || 'None received';
        }
    } catch (error) {
        console.error('Failed to check deep link:', error);
    }
}

// Initialize
document.addEventListener('DOMContentLoaded', function() {
    console.log('Hybrid Dashboard initialized');
    console.log('Connected to Go backend');

    // Initial update
    updateStats();

    // Update every second
    updateInterval = setInterval(updateStats, 1000);

    // Poll for deep links every 2 seconds
    setInterval(checkDeepLink, 2000);

    // Initial deep link check
    checkDeepLink();

    // Set up button handlers
    document.getElementById('callGoBtn').addEventListener('click', callGoFunction);
    document.getElementById('checkDeepLinkBtn').addEventListener('click', checkDeepLink);

    // Handle hash routing (for deep link navigation)
    window.addEventListener('hashchange', handleHashRoute);
    handleHashRoute(); // Check initial hash
});

// Handle hash-based routing from deep links
function handleHashRoute() {
    const hash = window.location.hash.slice(2); // Remove '#/'
    if (hash) {
        console.log('Route:', hash);
        // You could show/hide sections based on hash here
        // For now, just log it
    }
}

// Cleanup on page unload
window.addEventListener('beforeunload', function() {
    if (updateInterval) {
        clearInterval(updateInterval);
    }
});
