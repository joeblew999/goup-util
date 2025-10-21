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

// Initialize
document.addEventListener('DOMContentLoaded', function() {
    console.log('Hybrid Dashboard initialized');
    console.log('Connected to Go backend');
    
    // Initial update
    updateStats();
    
    // Update every second
    updateInterval = setInterval(updateStats, 1000);
    
    // Set up button handler
    document.getElementById('callGoBtn').addEventListener('click', callGoFunction);
});

// Cleanup on page unload
window.addEventListener('beforeunload', function() {
    if (updateInterval) {
        clearInterval(updateInterval);
    }
});
