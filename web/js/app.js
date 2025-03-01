// API endpoints
const API_BASE_URL = '/api/v1';
const CALCULATE_URL = `${API_BASE_URL}/calculate`;
const EXPRESSIONS_URL = `${API_BASE_URL}/expressions`;

// Helper function to show errors
function showError(message) {
    const errorElement = document.getElementById('errorMessage');
    errorElement.textContent = message;
    errorElement.classList.remove('hidden');
}

// Helper function to hide errors
function hideError() {
    const errorElement = document.getElementById('errorMessage');
    errorElement.classList.add('hidden');
}

// Submit expression to API
function submitExpression() {
    hideError();

    const expressionInput = document.getElementById('expression');
    const expression = expressionInput.value.trim();

    if (!expression) {
        showError('Expression cannot be empty');
        return;
    }

    fetch(CALCULATE_URL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ expression: expression }),
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to submit expression');
                });
            }
            return response.json();
        })
        .then(data => {
            // Show success message with expression ID
            document.getElementById('expressionId').textContent = data.id;
            document.getElementById('expressionLink').href = `/web/expression/${data.id}`;
            document.getElementById('result').classList.remove('hidden');
        })
        .catch(error => {
            showError(error.message);
        });
}

// Load all expressions
function loadExpressions() {
    hideError();

    const listElement = document.getElementById('expressionsList');
    listElement.innerHTML = '<li>Loading expressions...</li>';

    fetch(EXPRESSIONS_URL)
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to load expressions');
                });
            }
            return response.json();
        })
        .then(data => {
            if (!data.expressions || data.expressions.length === 0) {
                listElement.innerHTML = '<li>No expressions found</li>';
                return;
            }

            listElement.innerHTML = '';
            data.expressions.forEach(expr => {
                const li = document.createElement('li');
                li.className = 'expression-item';

                const status = document.createElement('span');
                status.className = `expression-status status-${expr.status}`;
                status.textContent = expr.status;

                li.innerHTML = `
                <p><strong>Expression:</strong> ${expr.expression}</p>
                <p><strong>Status:</strong> </p>
                <p>${expr.result !== undefined && expr.status === 'COMPLETE' ?
                    `<strong>Result:</strong> ${expr.result}` : ''}
                   ${expr.error ? `<strong>Error:</strong> <span class="error">${expr.error}</span>` : ''}</p>
                <p><a href="/web/expression/${expr.id}">View Details</a></p>
            `;

                li.querySelector('p:nth-child(2)').appendChild(status);
                listElement.appendChild(li);
            });
        })
        .catch(error => {
            listElement.innerHTML = '';
            showError(error.message);
        });
}

// Load expression detail
function loadExpressionDetail() {
    hideError();

    // Extract expression ID from URL
    const path = window.location.pathname;
    const expressionId = path.substring(path.lastIndexOf('/') + 1);

    if (!expressionId) {
        showError('Expression ID not provided');
        return;
    }

    fetch(`${EXPRESSIONS_URL}/${expressionId}`)
        .then(response => {
            if (!response.ok) {
                if (response.status === 404) {
                    throw new Error('Expression not found');
                }
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to load expression');
                });
            }
            return response.json();
        })
        .then(data => {
            const expr = data.expression;
            document.getElementById('expressionId').textContent = expr.id;
            document.getElementById('expressionValue').textContent = expr.expression;

            const statusElement = document.getElementById('expressionStatus');
            statusElement.textContent = expr.status;
            statusElement.className = `expression-status status-${expr.status}`;

            // Show result if available
            const resultContainer = document.querySelector('#expressionResult').parentNode;
            if (expr.result !== null && expr.result !== undefined) {
                document.getElementById('expressionResult').textContent = expr.result;
                resultContainer.classList.remove('hidden');
            } else {
                resultContainer.classList.add('hidden');
            }

            // Show error if available
            const errorContainer = document.querySelector('#expressionError').parentNode;
            if (expr.error) {
                document.getElementById('expressionError').textContent = expr.error;
                errorContainer.classList.remove('hidden');
            } else {
                errorContainer.classList.add('hidden');
            }
        })
        .catch(error => {
            showError(error.message);
        });
}

// Handle Enter key in expression input
document.addEventListener('DOMContentLoaded', function() {
    const expressionInput = document.getElementById('expression');
    if (expressionInput) {
        expressionInput.addEventListener('keypress', function(event) {
            if (event.key === 'Enter') {
                event.preventDefault();
                submitExpression();
            }
        });
    }
});