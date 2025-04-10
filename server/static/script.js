// Debug function to check what's being returned from the API
function logResponse(response) {
    console.log("Response:", response);
    return response;
}

// Load faculty sponsors on page load
function loadFacultySponsors() {
    fetch('/api/filters')
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log("Raw sponsors data:", data);
            const sponsorSelect = document.getElementById('sponsor-select');

            // Ensure we have an array of sponsors
            if (!data.facultySponsors || !Array.isArray(data.facultySponsors)) {
                console.error("Invalid sponsors data format:", data);
                return;
            }

            console.log(`Found ${data.facultySponsors.length} sponsors`);

            // Clear any existing options (except the default)
            sponsorSelect.innerHTML = '<option value="">Select a sponsor...</option>';

            // Add sponsors to dropdown
            data.facultySponsors.forEach(sponsor => {
                console.log(`Adding sponsor: ${sponsor}`);
                const option = document.createElement('option');
                option.value = sponsor;
                option.textContent = sponsor;
                sponsorSelect.appendChild(option);
            });
        })
        .catch(error => {
            console.error("Error fetching sponsors:", error);
            alert("Failed to load faculty sponsors. Please refresh the page.");
        });
}

// Process studies response 
function processStudiesResponse(event) {
    if (event.detail.target.id === 'study-select') {
        try {
            const rawData = event.detail.xhr.responseText;
            console.log("Raw studies response:", rawData);

            const data = JSON.parse(rawData);
            console.log("Parsed studies data:", data);
            const studySelect = document.getElementById('study-select');

            // Clear existing options
            studySelect.innerHTML = '<option value="">Select a study...</option>';

            // Ensure we have an array of studies
            if (!data.studies || !Array.isArray(data.studies)) {
                console.error("Invalid studies data format:", data);
                return;
            }

            console.log(`Found ${data.studies.length} studies`);

            // Add studies to dropdown
            data.studies.forEach(study => {
                console.log(`Adding study: ${study}`);
                const option = document.createElement('option');
                option.value = study;
                option.textContent = study;
                studySelect.appendChild(option);
            });
        } catch (e) {
            console.error("Error parsing studies response:", e);
        }
    }
}

// Handle sample data loading
function handleSampleDataLoaded(event) {
    if (event.detail.target.id === 'samples-container') {
        const sponsor = document.getElementById('sponsor-select').value;
        const study = document.getElementById('study-select').value;

        if (sponsor && study) {
            // Initialize pagination if we have a table
            initPagination();
            // Update chart
            updateChartWithFilters(sponsor, study);
        }
    }
}

// Initialize pagination for the samples table
function initPagination() {
    const tableBody = document.querySelector('#samples-container table tbody');
    if (!tableBody) return;

    const rows = Array.from(tableBody.querySelectorAll('tr'));
    const rowsPerPage = 25;
    const totalPages = Math.ceil(rows.length / rowsPerPage);

    if (totalPages <= 1) return; // No need for pagination

    // Create pagination controls
    const paginationInfo = document.createElement('div');
    paginationInfo.className = 'pagination-info';

    const paginationContainer = document.createElement('div');
    paginationContainer.className = 'pagination';

    // Add prev button
    const prevButton = document.createElement('button');
    prevButton.innerHTML = '&laquo;';
    prevButton.className = 'pagination-button';
    prevButton.disabled = true;
    paginationContainer.appendChild(prevButton);

    // Add page buttons
    for (let i = 1; i <= totalPages; i++) {
        const pageButton = document.createElement('button');
        pageButton.textContent = i;
        pageButton.className = i === 1 ? 'pagination-button active' : 'pagination-button';
        pageButton.dataset.page = i;
        paginationContainer.appendChild(pageButton);
    }

    // Add next button
    const nextButton = document.createElement('button');
    nextButton.innerHTML = '&raquo;';
    nextButton.className = 'pagination-button';
    paginationContainer.appendChild(nextButton);

    // Add pagination container after the table
    tableBody.parentElement.after(paginationInfo);
    paginationInfo.after(paginationContainer);

    // Set initial page
    showPage(1, rows, rowsPerPage, totalPages, paginationInfo);

    // Add event listeners for pagination buttons
    addPaginationEventListeners(paginationContainer, rows, rowsPerPage, totalPages, paginationInfo, prevButton, nextButton);
}

// Add event listeners to pagination buttons
function addPaginationEventListeners(paginationContainer, rows, rowsPerPage, totalPages, paginationInfo, prevButton, nextButton) {
    paginationContainer.addEventListener('click', function (e) {
        if (e.target.tagName !== 'BUTTON') return;

        const currentPage = parseInt(document.querySelector('.pagination-button.active').dataset.page) || 1;
        let targetPage = currentPage;

        if (e.target === prevButton && currentPage > 1) {
            targetPage = currentPage - 1;
        } else if (e.target === nextButton && currentPage < totalPages) {
            targetPage = currentPage + 1;
        } else if (e.target.dataset.page) {
            targetPage = parseInt(e.target.dataset.page);
        }

        if (targetPage !== currentPage) {
            // Update active button
            document.querySelectorAll('.pagination-button').forEach(btn => {
                btn.classList.remove('active');
                if (btn.dataset.page == targetPage) {
                    btn.classList.add('active');
                }
            });

            // Update prev/next button state
            prevButton.disabled = targetPage === 1;
            nextButton.disabled = targetPage === totalPages;

            // Show target page
            showPage(targetPage, rows, rowsPerPage, totalPages, paginationInfo);
        }
    });
}

// Show specified page of table rows
function showPage(pageNumber, rows, rowsPerPage, totalPages, infoElement) {
    const startIndex = (pageNumber - 1) * rowsPerPage;
    const endIndex = startIndex + rowsPerPage;

    // Update info text
    const totalRows = rows.length;
    infoElement.textContent = `Showing ${startIndex + 1} to ${Math.min(endIndex, totalRows)} of ${totalRows} entries`;

    // Show/hide rows
    rows.forEach((row, index) => {
        row.style.display = (index >= startIndex && index < endIndex) ? '' : 'none';
    });
}

// Update chart with filter values
function updateChartWithFilters(sponsor, study) {
    const params = new URLSearchParams();
    params.append('sponsor', sponsor);
    params.append('study', study);

    fetch('/api/chart?' + params.toString())
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            // Show the chart and hide the instruction box
            document.querySelector('#chart-container .instruction-box').classList.add('hidden');
            document.getElementById('timingChart').classList.remove('hidden');

            updateChart(data);
        })
        .catch(error => {
            console.error("Error fetching chart data:", error);
            alert("Failed to load chart data. Please try again.");
        });
}

// Create and update chart
let chart; // Global chart variable

function updateChart(data) {
    if (chart) {
        chart.destroy();
    }

    createChart(data);
}

function createChart(data) {
    // Check if we have data to display
    if (!data.labels || data.labels.length === 0) {
        console.log("No chart data available");
        return;
    }

    const ctx = document.getElementById('timingChart').getContext('2d');

    chart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: data.labels, // Using supplier names as labels
            datasets: [
                {
                    label: 'Library Time',
                    data: data.libraryTime,
                    backgroundColor: 'rgba(54, 162, 235, 0.7)',
                    borderColor: 'rgba(54, 162, 235, 1)',
                    borderWidth: 1
                },
                {
                    label: 'Sequencing Time',
                    data: data.sequencingTime,
                    backgroundColor: 'rgba(255, 99, 132, 0.7)',
                    borderColor: 'rgba(255, 99, 132, 1)',
                    borderWidth: 1
                }
            ]
        },
        options: {
            indexAxis: 'y',
            scales: {
                x: {
                    stacked: true,
                    title: {
                        display: true,
                        text: 'Days'
                    }
                },
                y: {
                    stacked: true,
                    title: {
                        display: true,
                        text: 'Samples'
                    }
                }
            },
            plugins: {
                title: {
                    display: true,
                    text: 'Library and Sequencing Processing Times'
                },
                tooltip: {
                    callbacks: {
                        footer: function (tooltipItems) {
                            const item = tooltipItems[0];
                            const index = item.dataIndex;
                            return 'Sample ID: ' + data.sampleIds[index];
                        }
                    }
                }
            }
        }
    });
}

// Setup event listeners
document.addEventListener('DOMContentLoaded', function () {
    // Load faculty sponsors on page load
    loadFacultySponsors();

    // Handle study select enabling/disabling
    document.getElementById('sponsor-select').addEventListener('change', function () {
        const studySelect = document.getElementById('study-select');
        const applyButton = document.getElementById('apply-filters');

        if (this.value) {
            studySelect.disabled = false;
        } else {
            studySelect.disabled = true;
            studySelect.innerHTML = '<option value="">Select a study...</option>';
            applyButton.disabled = true;
        }
    });

    // Enable apply button when both selections are made
    document.getElementById('study-select').addEventListener('change', function () {
        const sponsorSelect = document.getElementById('sponsor-select');
        const applyButton = document.getElementById('apply-filters');

        if (this.value && sponsorSelect.value) {
            applyButton.disabled = false;
        } else {
            applyButton.disabled = true;
        }
    });

    // Process HTMX responses
    document.body.addEventListener('htmx:afterSwap', function (event) {
        processStudiesResponse(event);
        handleSampleDataLoaded(event);
    });

    // Apply button handler - also updates the chart
    document.getElementById('apply-filters').addEventListener('click', function () {
        const sponsor = document.getElementById('sponsor-select').value;
        const study = document.getElementById('study-select').value;

        if (sponsor && study) {
            // HTMX will handle the sample table update
            // We manually trigger chart update here
            updateChartWithFilters(sponsor, study);
        }
    });
});