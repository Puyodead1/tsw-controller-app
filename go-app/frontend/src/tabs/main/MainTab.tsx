export const MainTab = () => {
  return (
    <div>
      <fieldset className="fieldset">
        <legend className="fieldset-legend">Select profile</legend>
        <select className="select w-full">
          <option disabled selected>
            Auto-detect
          </option>
        </select>
        <span className="label">Auto-detect only works for certain supported controllers</span>
      </fieldset>
    </div>
  );
}
