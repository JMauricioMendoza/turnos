import { Navbar, NavbarBrand, NavbarContent } from "@heroui/react";
import LogoICEO from "../assets/img/logoIceo.png";

function Layout({ children }) {
  return (
    <div className="min-h-screen w-screen bg-gray-100 flex flex-col gap-6 items-center">
      <Navbar isBordered>
        <NavbarBrand>
          <img src={LogoICEO} alt="Logo ICEO" className="h-12" />
        </NavbarBrand>

        <NavbarContent justify="end">
          <p className="font-bold text-xl text-institucional">
            Sistema de Turnos y Atención Ciudadana
          </p>
        </NavbarContent>
      </Navbar>
      {children}
    </div>
  );
}

export default Layout;
