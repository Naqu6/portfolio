document.querySelector(".copyright-year").textContent = (new Date).getFullYear();

function scrollToDiv(event, selector) {
	if (window.location.pathname == "/") {
		event.preventDefault();

		document.querySelector(selector).scrollIntoView({
			block: 'center',
			behavior: 'smooth'
		});
	}
}

document.querySelector(`a[href="/#portfolio-container"]`).addEventListener("click", (event) => {
	scrollToDiv(event, "#portfolio-container");
});

document.querySelector(`a[href="/#about"]`).addEventListener("click", (event) => {
	scrollToDiv(event, "#about");
});

document.querySelector(`a[href="/#resume"]`).addEventListener("click", (event) => {
	scrollToDiv(event, "#resume");
});

document.querySelector(".see-projects").addEventListener("click", () => {
	document.querySelector("#portfolio-container").scrollIntoView({
		block: 'center',
		behavior: 'smooth'
	});
});