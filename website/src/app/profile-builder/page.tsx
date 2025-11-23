import { Nav } from "@/_components/Nav";
import { ProfileEditor } from "@/_components/ProfileEditor.client";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "Profile Builder - TSW Controller App",
};

export default function ProfileBuilder() {
  return (
    <>
      <Nav />
      <ProfileEditor />
    </>
  );
}
