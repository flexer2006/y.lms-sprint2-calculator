// Common utility functions
function showError(message, elementId = 'errorMessage') {
    const errorElement = document.getElementById(elementId);
    errorElement.textContent = message;
    errorElement.classList.remove('hidden');
}

function hideError(elementId = 'errorMessage') {
    const errorElement = document.getElementById(elementId);
    errorElement.classList.add('hidden');
}

function getStatusClass(status) {
    return `status-${status}`;
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleString();
}

// Calculate page functions
function submitExpression() {
    const expression = document.getElementById('expression').value.trim();
    hideError();

    if (!expression) {
        showError('Please enter an expression');
        return;
    }

    fetch('/api/v1/calculate', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            expression: expression
        })
    })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to calculate expression');
                });
            }
            return response.json();
        })
        .then(data => {
            document.getElementById('result').classList.remove('hidden');
            document.getElementById('expressionId').textContent = data.id;
            document.getElementById('expressionLink').href = `/web/expressions/${data.id}`;
        })
        .catch(error => {
            showError(error.message);
        });
}

// Expression list page functions
function loadExpressions() {
    fetch('/api/v1/expressions')
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to load expressions');
                });
            }
            return response.json();
        })
        .then(data => {
            const expressionsList = document.getElementById('expressionsList');
            expressionsList.innerHTML = '';

            if (!data.expressions || data.expressions.length === 0) {
                expressionsList.innerHTML = '<li>No expressions found</li>';
                return;
            }

            data.expressions.forEach(expr => {
                const item = document.createElement('li');
                item.className = 'expression-item';

                const status = document.createElement('span');
                status.className = `expression-status ${getStatusClass(expr.status)}`;
                status.textContent = expr.status;

                item.innerHTML = `
                    <strong>ID:</strong> <a href="/web/expressions/${expr.id}">${expr.id}</a><br>
                    <strong>Expression:</strong> ${expr.expression || 'N/A'}<br>
                    <strong>Status:</strong> `;
                item.appendChild(status);

                if (expr.result !== undefined && expr.result !== null) {
                    item.innerHTML += `<br><strong>Result:</strong> ${expr.result}`;
                }

                if (expr.error) {
                    item.innerHTML += `<br><strong>Error:</strong> <span class="error">${expr.error}</span>`;
                }

                expressionsList.appendChild(item);
            });
        })
        .catch(error => {
            showError(error.message);
        });
}

// Expression detail page functions
function loadExpressionDetail() {
    const pathParts = window.location.pathname.split('/');
    const expressionId = pathParts[pathParts.length - 1];

    fetch(`/api/v1/expressions/${expressionId}`)
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to load expression details');
                });
            }
            return response.json();
        })
        .then(data => {
            const expr = data.expression;
            document.getElementById('expressionId').textContent = expr.id;
            document.getElementById('expressionValue').textContent = expr.expression || 'N/A';

            const statusElement = document.getElementById('expressionStatus');
            statusElement.textContent = expr.status;
            statusElement.className = `expression-status ${getStatusClass(expr.status)}`;

            const resultElement = document.getElementById('expressionResult');
            if (expr.result !== undefined && expr.result !== null) {
                resultElement.textContent = expr.result;
                resultElement.parentElement.classList.remove('hidden');
            } else {
                resultElement.parentElement.classList.add('hidden');
            }

            const errorElement = document.getElementById('expressionError');
            if (expr.error) {
                errorElement.textContent = expr.error;
                errorElement.parentElement.classList.remove('hidden');
            } else {
                errorElement.parentElement.classList.add('hidden');
            }
        })
        .catch(error => {
            showError(error.message);
        });
}