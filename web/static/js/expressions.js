async function loadExpressions() {
    let response = await fetch("/api/v1/expressions");
    let data = await response.json();
    let list = document.getElementById("expressions-list");

    list.innerHTML = "";
    data.expressions.forEach(exp => {
        let li = document.createElement("li");
        li.textContent = `ID: ${exp.id}, Статус: ${exp.status}, Результат: ${exp.result ?? "Ожидание"}`;
        list.appendChild(li);
    });
}

window.onload = loadExpressions;
