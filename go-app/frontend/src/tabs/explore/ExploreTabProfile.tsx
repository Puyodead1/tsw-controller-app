import { useState } from "react";
import {
  ImportSharedProfile,
  LoadConfiguration,
} from "../../../wailsjs/go/main/App";
import { main } from "../../../wailsjs/go/models";

type Props = {
  profile: main.Interop_SharedProfile;
};

export const ExploreTabProfile = ({ profile }: Props) => {
  const [downloading, setIsDownloading] = useState(false);
  const handleDownload = () => {
    setIsDownloading(true);
    ImportSharedProfile(profile)
      .then(() => LoadConfiguration().then(() => alert("Profile Downloaded")))
      .catch((err) => alert(String(err)))
      .finally(() => setIsDownloading(false));
  };

  return (
    <li className="list-row">
      <div className="list-col-grow">
        <div>{profile.Name}</div>
      </div>
      <div>
        <button className="btn btn-sm btn-primary" disabled={downloading} onClick={handleDownload}>
          {downloading && <span className="loading loading-spinner text-primary"></span>}
          Download Profile
        </button>
      </div>
    </li>
  );
};
