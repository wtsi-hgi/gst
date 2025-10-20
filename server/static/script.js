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

function loadSearchOptions() {
    var select = document.getElementById("search-select")
    
    var options = [
        "Sanger Sample ID",
        "Supplier Name",
        "Manifest Created",
        "Manifest Uploaded",
        "Labware Received",
        "Plate/Tube",
        "Order Made",
        "Library Start",
        "Library Complete",
        "Library Time",
        "Run ID",
        "Platform",
        "Pipeline",
        "Sequencing Run Start",
        "Sequencing QC Complete",
        "Sequencing Time",
        "QC Pass"
    ]

    options.forEach(option => {
                console.log(`Adding option: ${option}`);
                const opt = document.createElement('option');
                opt.value = option;
                opt.textContent = option;
                select.appendChild(opt);
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

class hotbarClass {
    constructor(paginationContainer, totalPages, visiblecount = 5) {
        this.paginationContainer = paginationContainer;
        this.totalPages = totalPages;
        this.visibleCount = visiblecount;
        this.selectedIndex = 1;
        this.leftmostIndex = 1;
        this.rightmostIndex = visiblecount;

        this.buttonMap = new Map();
        this._buildMap();
    }

    _buildMap(){
        // Add prev button
        const prevButton = document.createElement('button');
        prevButton.innerHTML = '&laquo;';
        prevButton.className = 'pagination-button';
        prevButton.disabled = true;
        this.paginationContainer.appendChild(prevButton);
        this.backButton = prevButton;

        // Add page select buttons
        for (let i = 1; i <= this.totalPages; i++) {
            const pageButton = document.createElement('button');
            pageButton.textContent = i;
            pageButton.className = i === 1 ? 'pagination-button active' : 'pagination-button';
            pageButton.dataset.page = i;

            // All buttons of index > visibleCount are invisible
            // TODO: consider the effect on page load performance here, since its
            // all done client side?
            if (i > this.visibleCount) {
                pageButton.hidden = true;
            }
            this.paginationContainer.appendChild(pageButton);

            this.buttonMap.set(i, pageButton)
        }

        // Add next button
        const nextButton = document.createElement('button');
        nextButton.innerHTML = '&raquo;';
        nextButton.className = 'pagination-button';
        this.paginationContainer.appendChild(nextButton);
        this.nextButton = nextButton;
    }

    getSelected() {
        return this.buttonMap.get(this.selectedIndex);
    }

    getBackButton() {
        return this.backButton;
    }

    getNextButton() {
        return this.nextButton;
    }

    selectIndex(index) {
        console.log("select index called for index", index)
        if (this.selectIndex != index) {
            // Update selected index active status
            const prev = this.getSelected();
            this.selectedIndex = index;
            prev.classList.remove('active');
            const next = this.getSelected();
            next.classList.add('active');

            // Update boundary indexes to maintain selected index near middle
            let middle = Math.floor(this.visibleCount / 2)
            let newLeftmost = this.selectedIndex - middle;
            newLeftmost = Math.max(0, newLeftmost);
            newLeftmost = Math.min(newLeftmost, this.buttonMap.size-this.visibleCount);
            let newRightmost = newLeftmost+this.visibleCount;

            // Update button visibility
            const difference = next.textContent - prev.textContent;
            
            if (difference > 0) { // Moving right
                for (let i = this.leftmostIndex; i < newLeftmost; i++) {
                    const left = this.buttonMap.get(i);
                    if (left) left.hidden = true;
                    const right = this.buttonMap.get(i+this.visibleCount);
                    if (right) right.hidden = false;
                }
            } else if (difference < 0) { // Moving left
                for(let i = this.rightmostIndex; i > newRightmost; i--) {
                    const right = this.buttonMap.get(i);
                    if (right) right.hidden = true;
                    const left = this.buttonMap.get(i-this.visibleCount);
                    if (left) left.hidden = false;
                }
            }
            
            this.leftmostIndex = newLeftmost;
            this.rightmostIndex = newRightmost;
        }
    }

    // TODO: re-add prev & next buttons, maybe change functionality? consider what
    // is appropriate.

    // TODO: add tests? is this possible since this is js based? does this imply
    // that this should be processed server side ? 

    // TODO: Add option to type page num?
}

var hotbar;

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

    hotbar = new hotbarClass(paginationContainer, totalPages, 10)

    // Add pagination container after the table
    tableBody.parentElement.after(paginationInfo);
    paginationInfo.after(paginationContainer);

    // Set initial page
    showPage(1, rows, rowsPerPage, totalPages, paginationInfo);

    // Add event listeners for pagination buttons
    prevButton = hotbar.getBackButton();
    nextButton = hotbar.getNextButton();
    addPaginationEventListeners(paginationContainer, rows, rowsPerPage, totalPages, paginationInfo, prevButton, nextButton);
}

// Add event listeners to pagination buttons
function addPaginationEventListeners(paginationContainer, rows, rowsPerPage, totalPages, paginationInfo, prevButton, nextButton) {
    paginationContainer.addEventListener('click', function (e) {
        console.log("clicked")
        if (e.target.tagName !== 'BUTTON') return;

        const currentPage = parseInt(document.querySelector('.pagination-button.active').dataset.page) || 1;
        let targetPage = currentPage;
        let nextButton = hotbar.getNextButton();
        let prevButton = hotbar.getBackButton();

        if (e.target === prevButton && currentPage > 1) {
            targetPage = currentPage - 1;
        } else if (e.target === nextButton && currentPage < totalPages) {
            targetPage = currentPage + 1;
        } else if (e.target.dataset.page) {
            targetPage = parseInt(e.target.dataset.page);
        }

        if (targetPage !== currentPage) {
            // Update active button
            hotbar.selectIndex(targetPage)

            // Update prev/next button state
            prevButton.disabled = targetPage === 1;
            nextButton.disabled = targetPage === totalPages;
            
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

// // Update chart to match search
// function updateChartWithSearch(searchText, searchCol) {
//     const params = new URLSearchParams();
//     params.append('text', searchText);
//     params.append('col', searchCol);

//     fetch('/api/chart?' + params.toString())
//         .then(response => {
//             if (!response.ok) {
//                 throw new Error('HTTP error ${response.status}')
//             }
//             return response.json();
//         })
//         .then(data => {

//         })
// }

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
            animation: false,
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
    loadSearchOptions();

    // Handle study select and search enabling/disabling
    document.getElementById('sponsor-select').addEventListener('change', function () {
        const studySelect = document.getElementById('study-select');
        const applyButton = document.getElementById('apply-filters');
        const searchMenu = document.getElementById('search-select');

        if (this.value) {
            studySelect.disabled = false;
        } else {
            studySelect.disabled = true;
            studySelect.innerHTML = '<option value="">Select a study...</option>';
            applyButton.disabled = true;
            searchMenu.innerHTML = '<option value=""> Select an option.. </option>';
            searchMenu.disabled = true;
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

    // Apply filter handler - also updates the chart & enables search
    document.getElementById('apply-filters').addEventListener('click', function () {
        const sponsor = document.getElementById('sponsor-select').value;
        const study = document.getElementById('study-select').value;
        const searchMenu = document.getElementById('search-select');
        const searchButton = document.getElementById('apply-search');

        if (sponsor && study) {
            // HTMX will handle the sample table update
            // We manually trigger chart update here
            updateChartWithFilters(sponsor, study);
            searchMenu.disabled = false;
            searchButton.disabled = false;
        }
    });

    // Apply search handler
    document.getElementById('apply-search').addEventListener('click', function () {
        const searchText = document.getElementById('query');
        const searchCol = document.getElementById('search-select');

        console.log("DEBUG: searchText: ", searchText.value, " searchCol: ", searchCol.value);
        if (searchText && searchCol) {
            // updateChartWithSearch(searchText, searchCol);
            // updateSampleTable(searchText, searchCol);
        }
    });
});