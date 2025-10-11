# **Guide: Configuring Steam Input to disable joystick / controller**

This guide explains how to use set-up and configure Steam Input to ensure the intercepted controllers are not being processed by Train Sim World 5/6.

---

## **Step 1: Open Steam and Access Game Properties**

1. Launch **Steam**.
2. Go to your **Library**.
3. Right-click the Train Sim World 5/6 and select **Properties**.

<a href='https://postimages.org/' target='_blank'><img src='https://i.postimg.cc/fLTNQ3CK/Screenshot-from-2025-09-27-22-11-04.png' border='0' alt='Screenshot-from-2025-09-27-22-11-04'/></a>

---

## **Step 2: Enable Steam Input**

1. In the **Properties** window, select the **Controller** tab.
2. Under **Override for [Game Name]**, select **Enable Steam Input**.
3. Close the properties window.

<a href='https://postimages.org/' target='_blank'><img src='https://i.postimg.cc/Mp2JK35h/Screenshot-from-2025-09-27-22-12-22.png' border='0' alt='Screenshot-from-2025-09-27-22-12-22'/></a>

> **Note:** This ensures Steam Input is always active for Train Sim World.

---

## **Step 3: Calibrate an Unrecognized Controller (if needed)**

If your controller is unrecognized by Steam (which is often the case for custom joysticks):

1. Go to **Steam → Settings → Controller**.
2. Connect your controller.
3. Select "Configure Controller" or "Calibration" for the controller you are trying to intercept.
4. In the calibration window, you will only need to configure the required "A" and "B" buttons; everything else can be skipped by re-hitting "A" repeatedly.

---

## **Step 4: Apply the Community Layout**

1. With your controller connected, click the **Controller Configuration** button in your game’s Library page.
2. Instead of browsing, use the direct link for the "Disabled Gamepad" layout:

**Train Sim World 5**
<a href="steam://controllerconfig/2967990/3576092503">steam://controllerconfig/2967990/3576092503</a>

**Train Sim World 6**
<a href="steam://controllerconfig/3656800/3576139582">steam://controllerconfig/3656800/3576139582</a>

<a href='https://postimages.org/' target='_blank'><img src='https://i.postimg.cc/5NZQ9JZv/Frame-15.jpg' border='0' alt='Disabled Gamepad'/></a>

3. Steam will open the layout and ask to **Apply Configuration**. Click **Apply**.

> **Note:** This particular template is designed to **disable all controller input for the game**, allowing another software to handle the controller.

---

## **Step 5: Verify the Layout**

1. Launch the game.
2. Test that the game ignores controller input
3. If the controller is still active in-game, double-check **Controller Configuration** and ensure the layout is applied.
