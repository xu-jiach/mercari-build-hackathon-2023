
import "./Header.css";
import logo from '../assets/logo.png'


export const Header: React.FC = () => {


  return (
    <>
      <header>
        <div>
          <img className="logo" src={logo} alt="logo" />
        </div>
      </header>
    </>
  );
}
