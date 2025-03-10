document.getElementById('submit').addEventListener('click', function() {
    const expression = document.getElementById('expression').value;
    if (expression) {
        fetch('http://localhost/api/v1/calculate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ expression: expression })
        })
        .then(response => {
            if (response.status === 201) {
                return response.json();
            } else {
                throw new Error('Ошибка при отправке выражения');
            }
        })
        .then(data => {
            console.log('Выражение добавлено с ID:', data.id);
            loadExpressions(); // Обновляем список выражений
        })
        .catch(error => {
            console.error(error);
        });
    }
});

function loadExpressions() {
    fetch('http://localhost/api/v1/expressions')
        .then(response => {
            if (response.status === 200) {
                return response.json();
            } else {
                throw new Error('Ошибка при загрузке выражений');
            }
        })
        .then(data => {
            const expressionsList = document.getElementById('expressions-list');
            expressionsList.innerHTML = ''; // Очищаем список
            data.expressions.forEach(expr => {
                const li = document.createElement('li');
                li.textContent = `ID: ${expr.id}, Статус: ${expr.status}, Результат: ${expr.result}`;
                expressionsList.appendChild(li);
            });
        })
        .catch(error => {
            console.error(error);
        });
}

// Загружаем список выражений при загрузке страницы
window.onload = loadExpressions;
