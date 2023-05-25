import { useCookies } from "react-cookie";
import "./Header.css";
import logo from '/workspace/mercari-build-hackathon-2023/frontend/simple-mercari-web/src/components/assets/logo.png'

export const Header: React.FC = () => {
  const [cookies, _, removeCookie] = useCookies(["userID", "token"]);

  const onLogout = (event: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    event.preventDefault();
    removeCookie("userID");
    removeCookie("token");
  };

  return (
    <>
      <header>
        <div>
          <img className="logo" src={logo} alt="logo" />
        </div>
        <div className="LogoutButtonContainer">
          <button onClick={onLogout} id="MerButton">
            Logout
          </button>
        </div>
      </header>
    </>
  );
}
