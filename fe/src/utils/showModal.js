import { createRoot } from "react-dom/client";
import ModalComp from "../components/ModalComp";

export function showModal({ variant, text, title = "", onAccept }) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);

  function closeModal() {
    root.unmount();
    container.remove();
  }

  const handleAccept = () => {
    if (onAccept) onAccept();
    closeModal();
  };

  root.render(
    <ModalComp
      isOpen
      onOpenChange={(isOpen) => {
        if (!isOpen) closeModal();
      }}
      variant={variant}
      text={text}
      title={title}
      onAccept={handleAccept}
    />,
  );
}
