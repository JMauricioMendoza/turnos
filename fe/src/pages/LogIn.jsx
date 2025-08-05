import { useState } from "react";
import {
  Card,
  CardHeader,
  CardBody,
  CardFooter,
  Input,
  Button,
  Form,
} from "@heroui/react";
import { IoMdEye, IoMdEyeOff } from "react-icons/io";
import { useNavigate } from "react-router-dom";
import Layout from "../components/Layout";
import { validateInput, removeSpaces } from "../utils/validators";
import sendData from "../utils/sendData";

function LogIn() {
  const [inputNombre, setInputNombre] = useState("");
  const [inputContrasena, setInputContrasena] = useState("");

  const [isVisible, setIsVisible] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  function onSuccess(data) {
    localStorage.setItem("token", data.datos);
    navigate("/inicio");
  }

  const onSubmit = (ev) => {
    sendData({
      ev,
      url: "/sesion/crear",
      method: "POST",
      dataToSend: {
        usuario: inputNombre,
        password: inputContrasena,
      },
      onSuccess,
      setIsLoading,
    });
  };

  return (
    <Layout>
      <div className="flex flex-col items-center gap-3">
        <h2 className="text-3xl font-bold">Iniciar sesión</h2>
        <p className="text-m text-zinc-700">
          Ingresa tus credenciales para acceder a tu cuenta
        </p>
      </div>
      <Form onSubmit={onSubmit}>
        <Card className="p-6 w-[450px]">
          <CardHeader>
            <h3 className="text-xl font-bold text-center w-full">Bienvenido</h3>
          </CardHeader>
          <CardBody className="flex flex-col gap-6">
            <Input
              label="Usuario"
              type="text"
              variant="bordered"
              isRequired
              value={inputNombre}
              onChange={(ev) => removeSpaces(ev, setInputNombre)}
            />
            <Input
              label="Contraseña"
              type={isVisible ? "text" : "password"}
              variant="bordered"
              isRequired
              value={inputContrasena}
              onChange={(ev) => removeSpaces(ev, setInputContrasena)}
              endContent={
                <button
                  className="self-center text-gray-700 text-xl"
                  type="button"
                  onClick={() => setIsVisible(!isVisible)}
                >
                  {isVisible ? <IoMdEyeOff /> : <IoMdEye />}
                </button>
              }
            />
          </CardBody>
          <CardFooter className="pt-6">
            <Button
              color="primary"
              type="submit"
              className="font-semibold text-base w-full"
              isDisabled={
                validateInput(inputNombre, { minLength: 6 }) ||
                validateInput(inputContrasena, { minLength: 6 })
              }
              isLoading={isLoading}
            >
              Iniciar sesión
            </Button>
          </CardFooter>
        </Card>
      </Form>
    </Layout>
  );
}

export default LogIn;
