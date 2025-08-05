import { showModal } from "./showModal";

async function sendData({
  ev = null,
  url,
  method,
  dataToSend = null,
  onSuccess = null,
  setIsLoading,
}) {
  if (ev?.preventDefault) ev.preventDefault();

  setIsLoading(true);

  const serverURL = process.env.REACT_APP_API_URL;

  const headers = {
    "Content-Type": "application/json",
  };

  try {
    const response = await fetch(`${serverURL}${url}`, {
      method,
      headers,
      ...(dataToSend ? { body: JSON.stringify(dataToSend) } : {}),
    });

    const data = await response.json().catch(() => ({}));
    setIsLoading(false);

    if (!response.ok) {
      switch (response.status) {
        case 400:
        case 401:
          showModal({
            variant: "warn",
            text: data.mensaje,
          });
          break;
        default:
          showModal({
            variant: "error",
            text: "No se ha podido completar la petición.",
          });
          break;
      }
      return;
    }

    if (onSuccess) {
      onSuccess(data);
    } else {
      showModal({
        variant: "ok",
        text: data.mensaje,
      });
    }
  } catch {
    setIsLoading(false);
    showModal({
      variant: "warn",
      text: "Ocurrió un error inesperado.",
    });
  }
}

export default sendData;
