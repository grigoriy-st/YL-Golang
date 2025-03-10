document.getElementById("calc-form").addEventListener("submit", async function (e) {
    e.preventDefault();
    let expression = document.getElementById("expression").value;
    let response = await fetch("/api/v1/calculate", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ expression }),
    });

    let result = await response.json();
    document.getElementById("response").textContent = response.ok ? `ID: ${result.id}` : "Ошибка!";
});
