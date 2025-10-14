import useSWR from "swr";
import {
  GetSharedProfiles,
  ImportSharedProfile,
  LoadConfiguration,
} from "../../../wailsjs/go/main/App";
import { main } from "../../../wailsjs/go/models";

export const ExploreTab = () => {
  const { data: sharedProfiles } = useSWR(
    "shared-profiles",
    () =>
      GetSharedProfiles().then((profiles) =>
        profiles.toSorted((a, b) => a.Name.localeCompare(b.Name)),
      ),
    { revalidateOnMount: true },
  );

  const handleDownload = (profile: main.Interop_SharedProfile) => {
    ImportSharedProfile(profile)
      .then(() => LoadConfiguration().then(() => alert("Profile Downloaded")))
      .catch((err) => alert(String(err)));
  };

  return (
    <div>
      <ul className="list bg-base-100 rounded-box shadow-md">
        {sharedProfiles?.map((profile) => (
          <li key={profile.Name} className="list-row">
            <div className="list-col-grow">
              <div>{profile.Name}</div>
            </div>
            <div>
              <button
                className="btn btn-sm btn-primary"
                onClick={() => handleDownload(profile)}
              >
                Download Profile
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};
