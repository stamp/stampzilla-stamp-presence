#include <SPI.h>           // SPI library
#include <SdFat.h>         // SDFat Library
#include <SdFatUtil.h>     // SDFat Util Library
#include <SFEMP3Shield.h>  // Mp3 Shield Library

SdFat sd; // Create object to handle SD functions

SFEMP3Shield MP3player; // Create Mp3 library object
// These variables are used in the MP3 initialization to set up
// some stereo options:
const uint8_t volume = 0; // MP3 Player volume 0=max, 255=lowest (off)
const uint16_t monoMode = 1;  // Mono setting 0=off, 3=max

/* Pin setup */
int lastTrigger = 0; // This variable keeps track of which tune is playing
byte status = 0;
byte prev = 0;

long lastDebounceTime = 0; 
long debounceDelay = 50;
int lastButtonState = HIGH;  
int buttonState; 

void initSD();
void initMP3Player();

void setup()
{
  pinMode(A0, INPUT);
  pinMode(A1, INPUT);
  pinMode(A2, INPUT);  
  pinMode(A3, INPUT);
  
  pinMode(A5, INPUT_PULLUP);

  Serial.begin(9600);

  initSD();  // Initialize the SD card
  initMP3Player(); // Initialize the MP3 Shield
}

// All the loop does is continuously step through the trigger
//  pins to see if one is pulled low. If it is, it'll stop any
//  currently playing track, and start playing a new one.
void loop()
{
  int reading = digitalRead(A5);
  if (reading != lastButtonState) {
    lastDebounceTime = millis();
  }

  if ((millis() - lastDebounceTime) > debounceDelay) {
    if (reading != buttonState) {
      buttonState = reading;
      if (buttonState == LOW) {
        Serial.println("<DOOR>");
        if (MP3player.isPlaying())
          MP3player.stopTrack();
    
        uint8_t result = MP3player.playMP3("trumpet.mp3");
  
        if (result != 0)  // playTrack() returns 0 on success
        {
          Serial.println("Failed");
        }
      }
    }
  }
  lastButtonState = reading;
  

  status = 0;
  if ((digitalRead(A0) == LOW) ) {
    status += 1; 
  }
  if ((digitalRead(A1) == LOW) ) {
    status += 2; 
  }
  if ((digitalRead(A2) == LOW) ) {
    status += 4; 
  }
  if ((digitalRead(A3) == LOW) ) {
    status += 8; 
  }
  
  if (status != prev) {
    Serial.print("<");
    for (unsigned int mask = 0x08; mask; mask >>= 1) {
       Serial.print(mask&status?'1':'0');
    }
    Serial.println(">");
    prev = status;    
  }
}

// initSD() initializes the SD card and checks for an error.
void initSD()
{
  //Initialize the SdCard.
  if(!sd.begin(SD_SEL, SPI_HALF_SPEED)) 
    sd.initErrorHalt();
  if(!sd.chdir("/")) 
    sd.errorHalt("sd.chdir");
}

// initMP3Player() sets up all of the initialization for the
// MP3 Player Shield. It runs the begin() function, checks
// for errors, applies a patch if found, and sets the volume/
// stero mode.
void initMP3Player()
{
  uint8_t result = MP3player.begin(); // init the mp3 player shield
  if(result != 0) // check result, see readme for error codes.
  {
    // Error checking can go here!
  }
  MP3player.setVolume(volume, volume);
  MP3player.setMonoMode(monoMode);
}
