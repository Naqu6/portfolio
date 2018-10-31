document.querySelector(".copyright-year").textContent = (new Date).getFullYear();

function scrollToDiv(event, selector) {
	if (window.location.pathname == "/") {
		event.preventDefault();

		document.querySelector(selector).scrollIntoView({
			block: 'start',
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

document.querySelector(`a[href="/#contact"]`).addEventListener("click", (event) => {
	scrollToDiv(event, "#contact");
});

if (window.location.pathname == "/") {
	document.querySelector(".see-projects").addEventListener("click", () => {
		document.querySelector("#portfolio-container").scrollIntoView({
			block: 'start',
			behavior: 'smooth'
		});
	});
}

function contact() {
	
	var textarea = document.querySelector("textarea");

	var successMessage = document.querySelector(".contact-success");

	successMessage.style.opacity = 0;

	if (!textarea.value) {
		successMessage.value = "Please Enter a Message";
		successMessage.style.opacity = 1;
	} else {
		successMessage.value = "Thank you for the message - I'll be in touch shortly";

		var xhr = new XMLHttpRequest();
		xhr.open('GET', `/contact?message=${encodeURI(textarea.value)}`);
		xhr.onload = function() {
		    if (xhr.status === 200) {
		        successMessage.style.opacity = 1;
		    }
		};
		xhr.send();

	} 
}