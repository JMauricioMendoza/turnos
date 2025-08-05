import {
  Modal,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalFooter,
  Button,
} from "@heroui/react";
import { MdError, MdWarning, MdCheckCircle } from "react-icons/md";

function ModalComp({
  isOpen,
  onOpenChange,
  variant,
  text,
  title = "",
  onAccept,
}) {
  const getHeaderIcon = () => {
    switch (variant) {
      case "error":
        return (
          <span className="text-danger text-2xl">
            <MdError />
          </span>
        );
      case "warn":
        return (
          <span className="text-warning text-2xl">
            <MdWarning />
          </span>
        );
      case "ok":
        return (
          <span className="text-success text-2xl">
            <MdCheckCircle />
          </span>
        );
      default:
        return null;
    }
  };

  const getTitle = () => {
    switch (variant) {
      case "error":
        return "Error";
      case "warn":
        return "Advertencia";
      case "ok":
        return "Proceso exitoso";
      default:
        return title;
    }
  };

  return (
    <Modal isOpen={isOpen} onOpenChange={onOpenChange}>
      <ModalContent>
        {(onClose) => (
          <>
            <ModalHeader>
              <div className="flex items-center gap-6">
                <p>{getTitle()}</p>
                {getHeaderIcon()}
              </div>
            </ModalHeader>
            <ModalBody>
              <p>{text}</p>
            </ModalBody>
            <ModalFooter>
              <Button
                color="primary"
                onPress={() => {
                  if (onAccept) onAccept();
                  onClose();
                }}
              >
                Aceptar
              </Button>
            </ModalFooter>
          </>
        )}
      </ModalContent>
    </Modal>
  );
}

export default ModalComp;
