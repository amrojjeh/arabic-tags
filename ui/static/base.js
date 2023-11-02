export function setOffline(offline = true) {
	if (offline) {
		document.querySelector("#offline-warning").classList.remove("hidden");
	} else {
		document.querySelector("#offline-warning").classList.add("hidden");
	}
}
