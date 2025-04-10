import handlers from "./build/worker.mjs";

const iframe = document.getElementById("worker-result");
const calculatorForm = document.getElementById("calculator");

calculatorForm.addEventListener("submit", async (e) => {
  e.preventDefault();

  const formData = new FormData(e.target);
  const a = formData.get("a");
  const b = formData.get("b");

  const req = new Request("/add", {
    method: "POST",
    body: JSON.stringify({ a: Number(a), b: Number(b) }),
  });
  const res = await handlers.fetch(req);
  const text = await res.text();

  iframe.srcdoc = text;
});
