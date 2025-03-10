async function loadExpression() {
    const id = window.location.pathname.split("/").pop();
    console.log(id)
    let response = await fetch(`/api/v1/expressions/${id}`);
    
    if (response.ok) {
        let data = await response.json();
        let expr = data.expression;
        document.getElementById("expression-details").innerHTML = `
            <p><strong>ID:</strong> ${expr.id}</p>
            <p><strong>Статус:</strong> ${expr.status}</p>
            <p><strong>Результат:</strong> ${expr.result ?? "Ожидание"}</p>
        `;
    } else {
        document.getElementById("expression-details").innerHTML = "<p>Выражение не найдено.</p>";
    }
}

window.onload = loadExpression;
